package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"
)

type Entry struct {
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Timestamp string            `json:"timestamp"`
	Tags      []string          `json:"tags"`
	Fields    map[string]string `json:"fields"`
	time_obj  time.Time
}

type Journal struct {
	Entries     []Entry `json:"days"`
	Last_loaded string  `json:"last_loaded"`
	Created     string  `json:"created"`
	Version     string  `json:"version"`
	Repo        string  `json:"-"`
	password    string
	path        string
	time_format string
}

// setters
func (j *Journal) setPassword(password string) {
	j.password = password
}

// getters
func (j *Journal) getCurrentVersion() (current_version string, e error) {
	// set a timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	// load url
	response, e := client.Get(j.Repo + "/releases/latest")
	if e != nil {
		return "", e
	}
	// get redirect
	final_url := strings.Split(response.Request.URL.String(), "/")
	if len(final_url) == 0 {
		return "", errors.New("unable to fetch new version")
	}
	// get the version tag
	current_version = final_url[len(final_url)-1]
	// remove the v
	current_version = current_version[1:]

	return current_version, e
}

func create_journal() (j Journal, e error) {
	// create journal path if it does not exist
	var journal_folder string
	// multi os support (hopefully)
	if runtime.GOOS == "linux" {
		journal_folder = "/var/lib/journal"
	} else if runtime.GOOS == "darwin" {
		// macOS, might need testing
		journal_folder = "~/Library/Preferences/journal"
	} else if runtime.GOOS == "windows" {
		journal_folder = os.Getenv("APPDATA") + "\\journal"
	}

	if _, e := os.Stat(journal_folder); os.IsNotExist(e) {
		e := os.Mkdir(journal_folder, 0777)
		fmt.Println(e)
		return Journal{}, errors.New("cannot create folder " + journal_folder)
	}

	j = Journal{
		Last_loaded: time.Now().Format(time.RFC3339),
		Version:     "1.1.0",
		Repo:        "https://github.com/lorossi/go-journal",
		time_format: "2006-01-02 15:04:05",
		path:        journal_folder + "/journal.json",
	}

	if strings.Contains(j.Version, "b") {
		j.path = "beta.json"
	}

	return j, nil
}

func read_from_file(path string) (file []byte, e error) {
	file, e = ioutil.ReadFile(path)

	if e != nil {
		return []byte("[]"), e
	}
	return file, e
}

func write_to_file(path string, bytes []byte) (e error) {
	e = ioutil.WriteFile(path, bytes, 0666)
	if e != nil {
		return errors.New("error while working with the file. cannot save")
	}
	return e
}

// package the variables into a new entry
func (j *Journal) createNewEntry(title, content string, tags []string, fields map[string]string, time_obj time.Time) (entry Entry) {
	var timestamp string
	// format the timestamp
	timestamp = time_obj.Format(j.time_format)
	// create the new entry
	entry = Entry{
		Title:     title,
		Content:   content,
		Tags:      tags,
		Fields:    fields,
		Timestamp: timestamp,
		time_obj:  time_obj,
	}

	return entry
}

// load entry from database
func (j *Journal) load() (e error) {
	var file []byte
	// try to open the file
	file, e = read_from_file(j.path)
	// if not available, create an empty one
	if e != nil {
		j.Created = time.Now().Format(time.RFC3339)
		ioutil.WriteFile(j.path, []byte("[]"), 0666)
		return nil
	}

	// parse JSON
	e = json.Unmarshal(file, &j)
	if e != nil {
		return errors.New("cannot parse database. Is it encrypted?")
	}
	// calculate the time for each entry
	for i := 0; i < len(j.Entries); i++ {
		j.Entries[i].time_obj, _ = time.Parse(j.time_format, j.Entries[i].Timestamp)
	}

	return nil
}

// load and decrypt database
func (j *Journal) decrypt() (e error) {
	var file []byte
	// try to open the file
	file, e = read_from_file(j.path)
	if e != nil {
		return errors.New("cannot open encrypted database")
	}

	key := []byte(j.password)

	c, e := aes.NewCipher(key)
	if e != nil {
		return errors.New("cannot create new cypher")
	}

	gcm, e := cipher.NewGCM(c)
	if e != nil {
		return errors.New("cannot create new GCM")
	}

	nonce_size := gcm.NonceSize()
	if len(file) < nonce_size {
		return errors.New("file is too short")
	}

	nonce, ciphertext := file[:nonce_size], file[nonce_size:]
	plaintext, e := gcm.Open(nil, nonce, ciphertext, nil)
	if e != nil {
		return errors.New("cannot decode file")
	}

	// parse JSON
	e = json.Unmarshal(plaintext, &j)
	if e != nil {
		return errors.New("cannot parse database")
	}
	// calculate the time for each entry
	for i := 0; i < len(j.Entries); i++ {
		j.Entries[i].time_obj, _ = time.Parse(j.time_format, j.Entries[i].Timestamp)
	}

	return nil
}

// save journal to database
func (j *Journal) save() (e error) {
	// Unmarshal data
	JSON_bytes, e := json.MarshalIndent(j, "", "  ")
	if e != nil {
		return errors.New("error while encoding data. cannot save")
	}
	// write to file
	write_to_file(j.path, JSON_bytes)
	return e
}

func (j *Journal) encrypt() (e error) {
	key := []byte(j.password)

	JSON_bytes, e := json.MarshalIndent(j, "", "  ")
	if e != nil {
		return errors.New("error while encoding data. cannot save")
	}

	c, e := aes.NewCipher(key)
	if e != nil {
		return errors.New("cannot create new cypher")
	}

	gcm, e := cipher.NewGCM(c)
	if e != nil {
		return errors.New("cannot create new GCM")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, e := io.ReadFull(rand.Reader, nonce); e != nil {
		return errors.New("cannot create new random sequence")
	}

	ciphertext := gcm.Seal(nonce, nonce, JSON_bytes, nil)
	// write to file
	e = write_to_file(j.path, ciphertext)
	return e
}

// create a new entry
func (j *Journal) createEntry(entry string) {
	// array of separators that end the title
	var delimiters = []string{".", "?", "!", "+", "@"}
	var current_delimiter string
	// title, content variables
	var title, content string
	// entry tags variable
	var tags []string
	// fields variable
	var fields map[string]string
	fields = make(map[string]string)
	// current datetime variable
	var new_date time.Time
	// new entry variable
	var new_entry Entry

	// find the submitted date and the entry without the (eventual) date
	content, new_date = parse_entry(entry, j.time_format)
	if content == "" {
		return
	}

	// find the delimiter between title and content
	current_delimiter = find_delimiter(content, delimiters)

	if current_delimiter == "" {
		// the title is the whole entry
		title = strings.TrimSpace(content)
		content = ""
	} else {
		split_entry := strings.Split(content, current_delimiter)
		// the title is the first part BEFORE the delimiter
		title = strings.TrimSpace(split_entry[0] + current_delimiter)
		content = strings.Replace(content, title, "", 1)
	}

	// now load the tags
	if strings.Contains(content, "+") {
		// we found one or more tags
		string_tags := strings.Split(content, "+")[1:]

		for _, s := range string_tags {
			tags = append(tags, strings.Split(s, " ")[0])
		}

		// now remove all tags from content
		for i := 0; i < len(tags); i++ {
			tags[i] = strings.TrimSpace(tags[i])
			content = strings.ReplaceAll(content, "+"+tags[i], "")
		}
	}

	// now load the fields
	if strings.Contains(content, "@") {
		// we found one or more fields
		string_fields := strings.Split(content, "@")[1:]

		for i := 0; i < len(string_fields); i++ {
			values := strings.Split(string_fields[i], "=")
			// if there arent' exactly 2 strings (key, value) separated by semicolon,
			// something is wrong
			if len(values) != 2 {
				print_error(errors.New("field '"+values[0]+"' provided in a wrong format"), 1)
			} else {
				fields[values[0]] = strings.TrimSpace(values[1])
			}

			// now remove all fields from content
			for i := 0; i < len(string_fields); i++ {
				content = strings.ReplaceAll(content, "@"+string_fields[i], "")
			}
		}
	}
	// remove leading / trailing spaces for a better Format
	content = strings.TrimSpace(content)
	// remove all multiple spaces
	content = remove_multiple_spaces(content)

	// finally, generate the new entry
	new_entry = j.createNewEntry(title, content, tags, fields, new_date)
	// append the entry to the entries array
	j.Entries = append(j.Entries, new_entry)
	// sort the entries array
	sort.Slice(j.Entries, func(i, k int) bool { return j.Entries[i].time_obj.Before(j.Entries[k].time_obj) })
}

func (j *Journal) removeEntry(timestamp string) (e error) {
	var remove_date time.Time
	var level int
	var clean_entries []Entry

	// get the date from the string
	remove_date, level = parse_day(timestamp)
	if level == -1 {
		return errors.New("date was not provided correctly")
	}

	// init an empty slice of entries
	clean_entries = make([]Entry, 0)
	for _, e := range j.Entries {
		// if the entry has the same date as the timestamp, don't append it
		// to the new slice of entries

		switch level {
		case 1:
			if !same_day(e.time_obj, remove_date) {
				clean_entries = append(clean_entries, e)
			}
		case 2:
			if !same_month(e.time_obj, remove_date) {
				clean_entries = append(clean_entries, e)
			}
		case 3:
			if !same_year(e.time_obj, remove_date) {
				clean_entries = append(clean_entries, e)
			}
		}
	}

	if len(clean_entries) == len(j.Entries) {
		// no entries were removed
		return errors.New("entry not found")
	} else {
		// replace the entries with a new slice
		j.Entries = clean_entries
		return nil
	}
}

func (j *Journal) showEntries(timestamp string) (entries []Entry, e error) {
	var get_date time.Time
	var level int

	entries = make([]Entry, 0)
	get_date, level = parse_day(timestamp)
	fmt.Println(get_date, level)

	// loop throught every entry and look for one with the desired day
	for _, e := range j.Entries {
		switch level {
		case 1:
			if same_day(e.time_obj, get_date) {
				entries = append(entries, e)
			}
		case 2:
			if same_month(e.time_obj, get_date) {
				entries = append(entries, e)
			}
		case 3:
			if same_year(e.time_obj, get_date) {
				entries = append(entries, e)
			}
		}
	}

	// if the loop has ended and none has been found, return an empty entry and
	// an error
	if len(entries) == 0 {
		return make([]Entry, 0), errors.New("no entries found")
	}

	// otherwise, return the entries
	return entries, nil
}

func (j *Journal) getAllEntries() ([]Entry, error) {
	if len(j.Entries) > 0 {
		return j.Entries, nil
	} else {
		// if there are no entries, return the empty slice and set an error
		return make([]Entry, 0), errors.New("no entries found")
	}
}

func (j *Journal) removeEntriesBetween(start_timestamp, end_timestamp string) (e error) {
	var start, end time.Time
	var clean_entries []Entry

	start, e = time.Parse("2006-01-02", start_timestamp)
	if e != nil {
		return errors.New("cannot parse start date")
	}
	end, e = time.Parse("2006-01-02", end_timestamp)
	if e != nil {
		return errors.New("cannot parse end date")
	}

	for _, entry := range j.Entries {
		if !date_between(entry.time_obj, start, end) {
			clean_entries = append(clean_entries, entry)
		}
	}

	if len(clean_entries) == len(j.Entries) {
		// no entries were removed
		return errors.New("entries not found")
	} else {
		// replace the entries with a new slice
		j.Entries = clean_entries
		return nil
	}
}

func (j *Journal) removeAllEntries() {
	j.Entries = make([]Entry, 0)
}

func (j *Journal) getEntriesBetween(start_timestamp, end_timestamp string) (entries []Entry, e error) {
	var start, end time.Time
	start, e = time.Parse("2006-01-02", start_timestamp)
	if e != nil {
		return make([]Entry, 0), errors.New("cannot parse start date")
	}
	end, e = time.Parse("2006-01-02", end_timestamp)
	if e != nil {
		return make([]Entry, 0), errors.New("cannot parse end date")
	}

	for _, entry := range j.Entries {
		if date_between(entry.time_obj, start, end) {
			entries = append(entries, entry)
		}
	}

	if len(entries) > 0 {
		return entries, nil
	} else {
		return make([]Entry, 0), errors.New("no entries found between those dates")
	}
}

func (j *Journal) searchKeywords(keywords []string) (entries []Entry, e error) {
	for _, entry := range j.Entries {
		for _, k := range keywords {
			if strings.Contains(entry.Title, k) || strings.Contains(entry.Content, k) {
				entries = append(entries, entry)
			}
		}
	}

	if len(entries) > 0 {
		return entries, nil
	} else {
		return make([]Entry, 0), errors.New("no entries found with the keyword")
	}
}

func (j *Journal) searchTags(tags []string) ([]Entry, error) {
	var entries []Entry
	for _, entry := range j.Entries {
		for _, entry_tag := range entry.Tags {
			for _, t := range tags {
				if entry_tag == t {
					entries = append(entries, entry)
					break
				}
			}
		}
	}

	if len(entries) > 0 {
		return entries, nil
	} else {
		return make([]Entry, 0), errors.New("no entries found with the tag")
	}
}

func (j *Journal) searchFields(keys []string) (entries []Entry, e error) {
	for _, entry := range j.Entries {
		for _, k := range keys {
			if _, ok := entry.Fields[k]; ok {
				entries = append(entries, entry)
			}
		}
	}

	if len(entries) > 0 {
		return entries, nil
	} else {
		return make([]Entry, 0), errors.New("no entries found with the field")
	}
}

func (j *Journal) getAllTags() (tags map[string]int, e error) {
	tags = make(map[string]int)
	for _, entry := range j.Entries {
		for _, tag := range entry.Tags {
			tags[tag] += 1
		}
	}

	if len(tags) > 0 {
		return tags, nil
	} else {
		return make(map[string]int), errors.New("no tags found")
	}
}

func (j *Journal) getAllFields() (fields []map[string]string, e error) {
	for _, entry := range j.Entries {
		if len(entry.Fields) > 0 {
			fields = append(fields, entry.Fields)
		}
	}

	if len(fields) > 0 {
		return fields, nil
	} else {
		return make([]map[string]string, 0), errors.New("no fields found")
	}
}

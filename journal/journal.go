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

// Entry contains a single entry in the journal
type Entry struct {
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Timestamp string            `json:"timestamp"`
	Tags      []string          `json:"tags"`
	Fields    map[string]string `json:"fields"`
	timeObj   time.Time
}

// Journal is the class containing the whole journal
type Journal struct {
	Entries          []Entry `json:"days"`
	LastLoaded       string  `json:"LastLoaded"`
	Created          string  `json:"created"`
	Version          string  `json:"version"`
	repo             string
	password         string
	folder, filename string
	timeFormat       string
}

//SetPassword -> sets new database password
func (j *Journal) SetPassword(password string) {
	j.password = password
}

//GetNewestVersion -> gets the current version from GitHub
func (j *Journal) GetNewestVersion() (newestVersion string, e error) {
	// set a timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	// load url
	response, e := client.Get(j.repo + "/releases/latest")
	if e != nil {
		return "", e
	}
	// get redirect
	finalURL := strings.Split(response.Request.URL.String(), "/")
	if len(finalURL) == 0 {
		return "", errors.New("unable to fetch new version")
	}
	// get the version tag
	newestVersion = finalURL[len(finalURL)-1]

	return newestVersion, e
}

// NewJournal returns an empty journal object
func NewJournal() (j Journal, e error) {
	// create journal path if it does not exist
	var journalFolder string
	// multi os support (hopefully)
	if runtime.GOOS == "linux" {
		journalFolder = "/var/lib/journal/"
	} else if runtime.GOOS == "darwin" {
		// macOS, might need testing
		journalFolder = "~/Library/Preferences/journal/"
	} else if runtime.GOOS == "windows" {
		journalFolder = os.Getenv("APPDATA") + "\\journal\\"
	}

	if _, e := os.Stat(journalFolder); os.IsNotExist(e) {
		e := os.Mkdir(journalFolder, 0777)
		fmt.Println(e)
		return Journal{}, errors.New("cannot create folder " + journalFolder)
	}

	j = Journal{
		LastLoaded: time.Now().Format(time.RFC3339),
		Version:    "1.1.3",
		repo:       "https://github.com/lorossi/go-journal",
		timeFormat: "2006-01-02 15:04:05",
		folder:     journalFolder,
		filename:   "journal.json",
	}

	if strings.Contains(j.Version, "b") {
		j.folder = ""
		j.filename = "beta.json"
	}

	return j, nil
}

func readFromFile(path string) (file []byte, e error) {
	file, e = ioutil.ReadFile(path)

	if e != nil {
		return []byte("[]"), e
	}
	return file, e
}

func writeToFile(path string, bytes []byte) (e error) {
	e = ioutil.WriteFile(path, bytes, 0666)
	if e != nil {
		return errors.New("error while working with the file. cannot save")
	}
	return e
}

// package the variables into a new entry
func (j *Journal) createNewEntry(title, content string, tags []string, fields map[string]string, timeObj time.Time) (entry Entry) {
	var timestamp string
	// format the timestamp
	timestamp = timeObj.Format(j.timeFormat)
	// create the new entry
	entry = Entry{
		Title:     title,
		Content:   content,
		Tags:      tags,
		Fields:    fields,
		Timestamp: timestamp,
		timeObj:   timeObj,
	}

	return entry
}

// load entry from database
func (j *Journal) load() (e error) {
	var file []byte
	// try to open the file
	file, e = readFromFile(j.folder + j.filename)

	// if not available, create an empty one
	// or, if JSON file is empty, just don't open it
	if e != nil || string(file) == "[]" {
		j.Created = time.Now().Format(time.RFC3339)
		ioutil.WriteFile(j.folder+j.filename, []byte("[]"), 0666)
		return nil
	}

	// parse JSON
	e = json.Unmarshal(file, &j)
	if e != nil {
		return errors.New("cannot parse database. Is it encrypted?")
	}

	// calculate the time for each entry
	for i := 0; i < len(j.Entries); i++ {
		j.Entries[i].timeObj, _ = time.Parse(j.timeFormat, j.Entries[i].Timestamp)
	}

	// update last loaded
	j.LastLoaded = time.Now().Format(time.RFC3339)

	return nil
}

func (j *Journal) setFilename(filename string) {
	j.filename = filename
}

// load and decrypt database
func (j *Journal) decrypt() (e error) {
	var file []byte
	// try to open the file
	file, e = readFromFile(j.folder + j.filename)
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

	nonceSize := gcm.NonceSize()
	if len(file) < nonceSize {
		return errors.New("file is too short")
	}

	nonce, ciphertext := file[:nonceSize], file[nonceSize:]
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
		j.Entries[i].timeObj, _ = time.Parse(j.timeFormat, j.Entries[i].Timestamp)
	}
	// update last loaded
	j.LastLoaded = time.Now().Format(time.RFC3339)

	return nil
}

// save journal to database
func (j *Journal) save() (e error) {
	// Unmarshal data
	JSONbytes, e := json.MarshalIndent(j, "", "  ")
	if e != nil {
		return errors.New("error while encoding data. cannot save")
	}
	// write to file
	writeToFile(j.folder+j.filename, JSONbytes)
	return e
}

func (j *Journal) encrypt() (e error) {
	key := []byte(j.password)

	JSONbytes, e := json.MarshalIndent(j, "", "  ")
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

	ciphertext := gcm.Seal(nonce, nonce, JSONbytes, nil)
	// write to file
	e = writeToFile(j.folder+j.filename, ciphertext)
	return e
}

// create a new entry
func (j *Journal) createEntry(entry string) {
	// array of separators that end the title
	var delimiters = []string{".", "?", "!", "+", "@"}
	var currentDelimiter string
	// title, content variables
	var title, content string
	// entry tags variable
	var tags []string
	// fields variable
	var fields map[string]string
	fields = make(map[string]string)
	// current datetime variable
	var newDate time.Time
	// new entry variable
	var newEntry Entry

	// find the submitted date and the entry without the (eventual) date
	content, newDate = parseEntry(entry)
	if content == "" {
		return
	}

	// find the delimiter between title and content
	currentDelimiter = findDelimiter(content, delimiters)

	if currentDelimiter == "" {
		// the title is the whole entry
		title = strings.TrimSpace(content)
		content = ""
	} else {
		splitEntry := strings.Split(content, currentDelimiter)
		// the title is the first part BEFORE the delimiter
		title = strings.TrimSpace(splitEntry[0] + currentDelimiter)
		content = strings.Replace(content, title, "", 1)
	}

	// now load the tags
	if strings.Contains(content, "+") {
		// we found one or more tags
		stringTags := strings.Split(content, "+")[1:]

		for _, s := range stringTags {
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
		sstringFields := strings.Split(content, "@")[1:]

		for i := 0; i < len(sstringFields); i++ {
			values := strings.Split(sstringFields[i], "=")
			// if there arent' exactly 2 strings (key, value) separated by semicolon,
			// something is wrong
			if len(values) != 2 {
				printError(errors.New("field '"+values[0]+"' provided in a wrong format"), 1)
			} else {
				fields[values[0]] = strings.TrimSpace(values[1])
			}

			// now remove all fields from content
			for i := 0; i < len(sstringFields); i++ {
				content = strings.ReplaceAll(content, "@"+sstringFields[i], "")
			}
		}
	}
	// remove leading / trailing spaces for a better Format
	content = strings.TrimSpace(content)
	// remove all multiple spaces
	content = removeMultipleSpaces(content)

	// finally, generate the new entry
	newEntry = j.createNewEntry(title, content, tags, fields, newDate)
	// append the entry to the entries array
	j.Entries = append(j.Entries, newEntry)
	// sort the entries array
	sort.Slice(j.Entries, func(i, k int) bool { return j.Entries[i].timeObj.Before(j.Entries[k].timeObj) })
}

func (j *Journal) removeEntry(timestamp string) (e error) {
	var removeDate time.Time
	var level int
	var cleanEntries []Entry

	// get the date from the string
	removeDate, level = parseDay(timestamp)
	if level == -1 {
		return errors.New("date was not provided correctly")
	}

	// init an empty slice of entries
	cleanEntries = make([]Entry, 0)
	for _, e := range j.Entries {
		// if the entry has the same date as the timestamp, don't append it
		// to the new slice of entries

		switch level {
		case 0:
			return errors.New("don't provide date with this flag")
		case 1:
			if !sameDay(e.timeObj, removeDate) {
				cleanEntries = append(cleanEntries, e)
			}
		case 2:
			if !sameMonth(e.timeObj, removeDate) {
				cleanEntries = append(cleanEntries, e)
			}
		case 3:
			if !sameYear(e.timeObj, removeDate) {
				cleanEntries = append(cleanEntries, e)
			}
		}
	}

	if len(cleanEntries) == len(j.Entries) {
		// no entries were removed
		return errors.New("entry not found")
	}
	// replace the entries with a new slice
	j.Entries = cleanEntries
	return nil

}

func (j *Journal) showEntries(timestamp string) (entries []Entry, e error) {
	var getDate time.Time
	var level int
	entries = make([]Entry, 0)
	getDate, level = parseDay(timestamp)

	// loop throught every entry and look for one with the desired day
	for _, e := range j.Entries {
		switch level {
		case 0:
			return entries, errors.New("don't provide date with this flag")
		case 1:
			if sameDay(e.timeObj, getDate) {
				entries = append(entries, e)
			}
		case 2:
			if sameMonth(e.timeObj, getDate) {
				entries = append(entries, e)
			}
		case 3:
			if sameYear(e.timeObj, getDate) {
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
	}

	// if there are no entries, return the empty slice and set an error
	return make([]Entry, 0), errors.New("no entries found")

}

func (j *Journal) removeEntriesBetween(startTimestamp, endTimestamp string) (e error) {
	var start, end time.Time
	var cleanEntries []Entry

	start, e = time.Parse("2006-01-02", startTimestamp)
	if e != nil {
		return errors.New("cannot parse start date")
	}
	end, e = time.Parse("2006-01-02", endTimestamp)
	if e != nil {
		return errors.New("cannot parse end date")
	}

	for _, entry := range j.Entries {
		if !dateBetween(entry.timeObj, start, end) {
			cleanEntries = append(cleanEntries, entry)
		}
	}

	if len(cleanEntries) == len(j.Entries) {
		// no entries were removed
		return errors.New("entries not found")
	}
	// replace the entries with a new slice
	j.Entries = cleanEntries
	return nil

}

func (j *Journal) removeAllEntries() {
	j.Entries = make([]Entry, 0)
}

func (j *Journal) getEntriesBetween(startTimestamp, endTimestamp string) (entries []Entry, e error) {
	var start, end time.Time
	start, e = time.Parse("2006-01-02", startTimestamp)
	if e != nil {
		return make([]Entry, 0), errors.New("cannot parse start date")
	}
	end, e = time.Parse("2006-01-02", endTimestamp)
	if e != nil {
		return make([]Entry, 0), errors.New("cannot parse end date")
	}

	for _, entry := range j.Entries {
		if dateBetween(entry.timeObj, start, end) {
			entries = append(entries, entry)
		}
	}

	if len(entries) > 0 {
		return entries, nil
	}
	return make([]Entry, 0), errors.New("no entries found between those dates")

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
	}
	return make([]Entry, 0), errors.New("no entries found with the keyword")

}

func (j *Journal) searchTags(tags []string) ([]Entry, error) {
	var entries []Entry
	for _, entry := range j.Entries {
		for _, entryTag := range entry.Tags {
			for _, t := range tags {
				if entryTag == t {
					entries = append(entries, entry)
					break
				}
			}
		}
	}

	if len(entries) > 0 {
		return entries, nil
	}
	return make([]Entry, 0), errors.New("no entries found with the tag")

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
	}
	return make([]Entry, 0), errors.New("no entries found with the field")

}

func (j *Journal) getAllTags() (tags map[string]int, e error) {
	tags = make(map[string]int)
	for _, entry := range j.Entries {
		for _, tag := range entry.Tags {
			tags[tag]++
		}
	}

	if len(tags) > 0 {
		return tags, nil
	}
	return make(map[string]int), errors.New("no tags found")

}

func (j *Journal) getAllFields() (fields []map[string]string, e error) {
	for _, entry := range j.Entries {
		if len(entry.Fields) > 0 {
			fields = append(fields, entry.Fields)
		}
	}

	if len(fields) > 0 {
		return fields, nil
	}
	return make([]map[string]string, 0), errors.New("no fields found")

}

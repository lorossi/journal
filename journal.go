package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	"time"
)

type Fields struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Entry struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Timestamp string   `json:"timestamp"`
	Tags      []string `json:"tags"`
	Fields    []Fields `json:"fields"`
	time_obj  time.Time
}

type Journal struct {
	Entries     []Entry `json:"days"`
	Last_loaded string  `json:"last_loaded"`
	Created     string  `json:"created"`
	Version     string  `json:"version"`
	path        string
	time_format string
}

func crate_journal() Journal {
	j := Journal{
		path:        "database.json",
		time_format: "2006-01-02",
		Last_loaded: time.Now().Format(time.RFC3339),
		Version:     "0.0.1",
	}

	return j
}

// package the variables into a new entru
func (j *Journal) createNewEntry(title, content string, tags []string, time_obj time.Time) Entry {
	var timestamp string
	// format the timestamp
	timestamp = time_obj.Format(j.time_format)
	// create the new entry
	entry := Entry{
		Title:     title,
		Content:   content,
		Tags:      tags,
		Timestamp: timestamp,
		time_obj:  time_obj,
	}

	return entry
}

// load entry from database
func (j *Journal) load() {
	// try to open the file
	file, e := ioutil.ReadFile(j.path)

	// if not available, create an empty one
	if e != nil {
		j.Created = time.Now().Format(time.RFC3339)
		ioutil.WriteFile(j.path, []byte("[]"), 0666)
		return
	}

	e = json.Unmarshal([]byte(file), &j)
	if e != nil {
		return
	}

	for i := 0; i < len(j.Entries); i++ {
		j.Entries[i].time_obj, _ = time.Parse(j.time_format, j.Entries[i].Timestamp)
	}
}

// save journal to database
func (j *Journal) save() {
	// Unmarshal data
	JSON_bytes, _ := json.MarshalIndent(j, "", "  ")
	// write to file
	_ = ioutil.WriteFile(j.path, JSON_bytes, 0666)
}

// create a new entry
func (j *Journal) createEntry(entry string) {
	// array of separators that end the title
	var delimiters = []string{".", ",", "?", "!"}
	var current_delimiter string
	// title, content variables
	var title, content string
	// entry split by words, tags variables
	var tags []string
	// current datetime variable
	var new_date time.Time
	// new entry variable
	var new_entry Entry

	// find the submitted date and the entry without the (eventual) date
	content, new_date = parse_day_entry(entry, j.time_format)
	// find the delimiter between title and content
	current_delimiter = find_delimiter(entry, delimiters)

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
		tags = strings.Split(content, "+")[1:]

		// now remove all tags from content
		for _, tag := range tags {
			content = strings.ReplaceAll(content, "+"+tag, "")
		}
	}

	// remove leading / trailing spaces for a better Format
	content = strings.TrimSpace(content)

	// finally, generate the new entry
	new_entry = j.createNewEntry(title, content, tags, new_date)
	// append the entry to the entries array
	j.Entries = append(j.Entries, new_entry)
}

func (j *Journal) removeEntry(timestamp string) error {
	var remove_date time.Time
	var clean_entries []Entry

	// get the date from the string
	remove_date = parse_day(timestamp, j.time_format)
	// init an empty slice of entries
	clean_entries = clean_entries[:0]
	for _, e := range j.Entries {
		// if the entry has the same date as the timestamp, don't append it
		// to the new slice of entries
		if !same_date(e.time_obj, remove_date) {
			clean_entries = append(clean_entries, e)
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

func (j *Journal) getEntry(timestamp string) (Entry, error) {
	var get_date time.Time
	get_date = parse_day(timestamp, j.time_format)

	// loop throught every entry and look for one with the desired day
	for _, e := range j.Entries {
		if same_date(e.time_obj, get_date) {
			return e, nil
		}
	}

	// if the loop has ended and none has been found, return an empty entry and
	// an error
	return Entry{}, errors.New("entry not found")
}

func (j *Journal) getAllEntries() ([]Entry, error) {
	if len(j.Entries) > 0 {
		return j.Entries, nil
	} else {
		// if there are no entries, return the empty slice and set an erro
		return j.Entries, errors.New("no entries found")
	}
}

func (j *Journal) searchKeywords(keywords []string) ([]Entry, error) {
	var entries []Entry
	for _, e := range j.Entries {
		for _, k := range keywords {
			if strings.Contains(e.Title, k) || strings.Contains(e.Content, k) {
				entries = append(entries, e)
			}
		}
	}

	if len(entries) > 0 {
		return entries, nil
	} else {
		return entries, errors.New("no entries found with the keyword")
	}
}

func (j *Journal) searchTags(tags []string) ([]Entry, error) {
	var entries []Entry
	for _, e := range j.Entries {
		for _, entry_tag := range e.Tags {
			for _, t := range tags {
				if entry_tag == t {
					entries = append(entries, e)
					break
				}
			}
		}
	}

	if len(entries) > 0 {
		return entries, nil
	} else {
		return entries, errors.New("no entries found with the tag")
	}
}

func (j *Journal) getMonth(month string) ([]Entry, error) {
	var entries []Entry
	var time_obj time.Time
	var e error

	time_obj, e = time.Parse("2006-01", month)
	if e != nil {
		return j.Entries, errors.New("date was not provided correctly")
	}

	for _, e := range j.Entries {
		if same_month(e.time_obj, time_obj) {
			entries = append(entries, e)
		}
	}

	if len(entries) > 0 {
		return entries, nil
	} else {
		return entries, errors.New("no entries found within this month")
	}
}

func (j *Journal) getYear(year string) ([]Entry, error) {
	var entries []Entry
	var time_obj time.Time
	var e error

	time_obj, e = time.Parse("2006", year)
	if e != nil {
		return j.Entries, errors.New("date was not provided correctly")
	}

	for _, e := range j.Entries {
		if same_year(e.time_obj, time_obj) {
			entries = append(entries, e)
		}
	}

	if len(entries) > 0 {
		return entries, nil
	} else {
		return entries, errors.New("no entries found within this year")
	}
}

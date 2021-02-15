package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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
	path        string
	time_format string
}

func crate_journal() (j Journal) {
	j = Journal{
		path:        "database.json",
		time_format: "2006-01-02",
		Last_loaded: time.Now().Format(time.RFC3339),
		Version:     "0.0.1",
	}

	return j
}

// package the variables into a new entru
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
	var delimiters = []string{".", ",", "?", "!", "+", "@"}
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

	// finally, generate the new entry
	new_entry = j.createNewEntry(title, content, tags, fields, new_date)
	// parse the timestamp and check if an entry for this day already exists
	timestamp := time.Now().Format("2006-01-02")
	for _, e := range j.Entries {
		if e.Timestamp == timestamp {
			print_error(errors.New("entry already found for this day. Aborting"), 1)
			return
		}
	}
	// append the entry to the entries array
	j.Entries = append(j.Entries, new_entry)
}

func (j *Journal) removeEntry(timestamp string) (e error) {
	var remove_date time.Time
	var level int8
	var clean_entries []Entry

	// get the date from the string
	remove_date, level = parse_day(timestamp)
	if level == -1 {
		return errors.New("date was not provided correctly")
	}
	// init an empty slice of entries
	clean_entries = clean_entries[:0]
	for _, e := range j.Entries {
		// if the entry has the same date as the timestamp, don't append it
		// to the new slice of entries

		switch level {
		case 0:
			if !same_day(e.time_obj, remove_date) {
				clean_entries = append(clean_entries, e)
			}
		case 1:
			if !same_month(e.time_obj, remove_date) {
				clean_entries = append(clean_entries, e)
			}
		case 2:
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

func (j *Journal) viewEntry(timestamp string) (entry Entry, e error) {
	var get_date time.Time
	var level int8

	get_date, level = parse_day(timestamp)

	// loop throught every entry and look for one with the desired day
	for _, e := range j.Entries {
		switch level {
		case 0:
			if same_day(e.time_obj, get_date) {
				return e, nil
			}
		case 1:
			if same_month(e.time_obj, get_date) {
				return e, nil
			}
		case 2:
			if same_year(e.time_obj, get_date) {
				return e, nil
			}
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
		// if there are no entries, return the empty slice and set an error
		return make([]Entry, 0), errors.New("no entries found")
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
		return make([]Entry, 0), errors.New("no entries found with the keyword")
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
		return make([]Entry, 0), errors.New("no entries found with the tag")
	}
}

func (j *Journal) searchFields(keys []string) (entries []Entry, e error) {
	for _, e := range j.Entries {
		for _, k := range keys {
			if _, ok := e.Fields[k]; ok {
				entries = append(entries, e)
			}
		}
	}

	if len(entries) > 0 {
		return entries, nil
	} else {
		return make([]Entry, 0), errors.New("no entries found with the field")
	}
}

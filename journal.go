package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

type Fields struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Entry struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Timestamp string    `json:"timestamp"`
	Tags      []string  `json:"tags"`
	Fields    []Fields  `json:"fields"`
	Time_obj  time.Time `json:"-"`
}

type Journal struct {
	Entries     []Entry  `json:"days"`
	Args        []string `json:"-"`
	path        string
	time_format string
}

func crate_journal() Journal {
	j := Journal{
		path:        "database.json",
		time_format: "2006-01-02",
		Args:        []string{"--add", "--remove", "--view", "--search", "--edit"},
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
		Time_obj:  time_obj,
	}

	return entry
}

// append the entry to the entries array
func (j *Journal) addNewEntry(new_entry Entry) {
	j.Entries = append(j.Entries, new_entry)
}

// load entry from database
func (j *Journal) load() {
	// try to open the file
	file, e := ioutil.ReadFile(j.path)

	// if not available, create an empty one
	if e != nil {
		ioutil.WriteFile(j.path, []byte("[]"), 0666)
		return
	}

	e = json.Unmarshal([]byte(file), &j)
	if e != nil {
		return
	}

	for i := 0; i < len(j.Entries); i++ {
		j.Entries[i].Time_obj, _ = time.Parse(j.time_format, j.Entries[i].Timestamp)
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
	var words, tags []string
	// current datetime variable
	var time_obj time.Time
	// new entry variable
	var new_entry Entry

	// check if the first word is either today or yesterday
	switch words = strings.Split(entry, " "); strings.ToLower(words[0]) {
	case "yesterday:":
		time_obj = time.Now().AddDate(0, 0, -1)
		content = strings.Join(words[1:], " ")
	case "today:":
		content = strings.Join(words[1:], " ")
		time_obj = time.Now()
	default:
		// try to parse time
		var e error
		time_obj, e = time.Parse(j.time_format, strings.TrimSpace(words[0]))
		if e == nil {
			content = strings.Join(words[1:], " ")
		} else {
			// failed, time is now
			content = entry
			time_obj = time.Now()
		}
	}

	// find the delimiter between title and content
	current_delimiter = find_delimiter(entry, delimiters)

	if current_delimiter == "" {
		// the title is the whole entry
		title = strings.TrimSpace(content)
	} else {
		split_entry := strings.Split(content, current_delimiter)
		// the title is the first part BEFORE the delimiter
		title = strings.TrimSpace(split_entry[0] + current_delimiter)
	}

	// now load the tags
	if strings.Contains(content, "+") {
		// we found one or more tags
		tags = strings.Split(content, "+")[1:]
	}

	// finally, generate and add the new entry
	new_entry = j.createNewEntry(title, content, tags, time_obj)
	j.addNewEntry(new_entry)
}

func (j *Journal) showDay(timestamp string) {
	for _, e := range j.Entries {
		if e.Timestamp == timestamp {
			fmt.Println("Date: ", e.Timestamp)
			fmt.Println("Title: ", e.Title)
			fmt.Println("Content: ", e.Content)
			fmt.Println("Tags:")
			for _, t := range e.Tags {
				fmt.Println("\t", t)
			}
			fmt.Println("Fields:")
			for _, f := range e.Fields {
				fmt.Println("\t", f.Key, ": ", f.Value)
			}
			return
		}
	}
}

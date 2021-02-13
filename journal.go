package main

import (
  "time"
  "io/ioutil"
  "encoding/json"
  "fmt"
  "strings"
)

type Fields struct {
  Key string `json:"key"`
  Value string `json:"value"`
}

type Entry struct {
  Title string `json:"title"`
  Content string `json:"content"`
  Timestamp string `json:"timestamp"`
  Tags []string `json:"tags"`
  Fields  []Fields `json:"fields"`
  Time_obj time.Time
}

type Journal struct {
  Entries []Entry `json:"days"`
  Path string
}

// package the variables into a new entru
func create_new_entry(title, content string, time_obj time.Time) Entry {
  var timestamp string
  // format the timestamp
  timestamp = time_obj.Format("2006-01-02, 03:04")
  // create the new entry
  entry := Entry {
    Title: title,
    Content: content,
    Timestamp: timestamp,
  }

  return entry
}

// append the entry to the entries array
func (j * Journal) addNewEntry(new_entry Entry) {
  j.Entries = append(j.Entries, new_entry)
}

// load entry from database
func (j * Journal) load() {
  // try to open the file
  file, e := ioutil.ReadFile(j.Path)

  // if not available, create an empty one
  if e != nil {
    ioutil.WriteFile(j.Path, []byte("[]"), 0666)
    return
  }
  _ = json.Unmarshal([]byte(file), &j)

  for i := 0; i < len(j.Entries); i++ {
    j.Entries[i].Time_obj, _ = time.Parse("2006-01-02, 03:04", j.Entries[i].Timestamp)
  }
}

// save journal to database
func (j * Journal) save() {
  // Unmarshal data
  JSON_bytes, _ := json.MarshalIndent(j.Entries, "", "  ")
  // write to file
  _ = ioutil.WriteFile(j.Path, JSON_bytes, 0666);
}

// create a new entry
func (j * Journal) createEntry(entry string) {
  // array of separators that end the title
  var delimiters = []string{".", ",", "?", "!"}
  var current_delimiter string
  // title, content variables
  var title, content string
  // current datetime variable
  var time_obj time.Time
  // new entry variable
  var new_entry Entry

  // find the delimiter between title and content
  current_delimiter = find_delimiter(entry, delimiters)

  if current_delimiter == "" {
    // the title is the whole entry
    title = strings.TrimSpace(entry)
  } else {
    split_entry := strings.Split(entry, current_delimiter)
    // the title is the first part BEFORE the delimiter
    title = strings.TrimSpace(split_entry[0] + current_delimiter)
  }
  // the content is the whole entry
  content = entry
  // load current time
  time_obj = time.Now()
  // finally, generate and add the new entry
  new_entry = create_new_entry(title, content, time_obj)
  j.addNewEntry(new_entry)
}

func (j * Journal) showDay(Timestamp string) {
  for _, e := range j.Entries {
    if e.Timestamp == Timestamp {
      fmt.Println("Date: ", e.Time_obj)
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
      break
    }
  }
}

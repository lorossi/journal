package main

import (
  "time"
  "encoding/json"
	"fmt"
	"io/ioutil"
  "flag"
  "bufio"
  "os"
  "strings"
)

type Fields struct {
  Key string `json:"key"`
  Value string `json:"value"`
}

type Day struct {
  Title string  `json:"title"`
  Content string `json:"content"`
  Timestamp string `json:"timestamp"`
  Tags []string `json:"tags"`
  Fields  []Fields `json:"fields"`
  Time_obj time.Time
}

type Journal struct {
  Days []Day `json:"days"`
}

func (j * Journal) load(path string) {
  // try to open the file
  file, e := ioutil.ReadFile(path)

  // if not available, create an empty one
  if e != nil {
    ioutil.WriteFile(path, []byte("[]"), 0666)
    return
  }
  _ = json.Unmarshal([]byte(file), &j)

  for i := 0; i < len(j.Days); i++ {
    j.Days[i].Time_obj, _ = time.Parse("2006-01-02", j.Days[i].Timestamp)
  }
}

func (j * Journal) addEntry() {
  reader := bufio.NewReader(os.Stdin)
  fmt.Println("New entry:")
  entry, _ := reader.ReadString('\n')

  var separator string = "."
  var split_entry []string = strings.Split(entry, separator)
  var title, content string
  title = split_entry[0] + separator

  if len(split_entry) == 1 {
    content = ""
  } else {
    content = strings.Join(split_entry[1:], separator)
  }

  // continue...

}

func (j * Journal) showDay(Timestamp string) {
  for _, d := range j.Days {
    if d.Timestamp == Timestamp {
      fmt.Println("Date: ", d.Time_obj)
      fmt.Println("Title: ", d.Title)
      fmt.Println("Content: ", d.Content)
      fmt.Println("Tags:")
      for _, t := range d.Tags {
        fmt.Println("\t", t)
      }
      fmt.Println("Fields:")
      for _, f := range d.Fields {
        fmt.Println("\t", f.Key, ": ", f.Value)
      }
      break
    }
  }
}

func main() {
  // vreate empty Journal
  j := Journal{}
  // load from database
  j.load("database.json")

  // parse all the args
  flag.Parse()
  // no args were found, time to add a new entry
  if flag.NArg() == 0 {
    j.addEntry()
  }
}

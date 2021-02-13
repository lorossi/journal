package main

import (
  "fmt"
  "strings"
  "os"
)

func main() {
  // vreate empty Journal
  j := Journal{Path: "database.json"}
  // load from database
  j.load()
  // declare new entry that will be parsed from console
  var entry string;

  // loop thought all arguments
  switch args := len(os.Args); {
  case args == 1:
    // user has provided no args. Prompt him to say something
    fmt.Println("Provide some args. Available args:")
    return

  case args > 1:
    // the user has provided some args. Join all together and add these as new entry
    entry = strings.Join(os.Args[1:], " ")
  }

  // add entry to Journal
  j.addEntry(entry)
  // save journal
  j.save()
}

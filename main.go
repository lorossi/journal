package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// vreate empty Journal
	j := crate_journal()
	// load from database
	j.load()
	// declare new entry that will be parsed from console
	var entry string

	// loop thought all arguments
	switch args := len(os.Args); {
	case args == 1:
		// user has provided no args. Prompt him to say something
		fmt.Println("Write your entry or provide some args. Available args:")
		fmt.Println(strings.Join(j.Args, "\n"))
		return

	case args > 1:
		// detect the passed arg
		switch param := strings.ToLower(os.Args[1]); param {
		case "--remove":
			fmt.Println("REMOVE")
		case "--view":
			fmt.Println("VIEW")
		case "--search":
			fmt.Println("SEARCH")
		case "--add":
			// the user has provided the "--add" flag. Time to add an entry.
			entry = strings.Join(os.Args[2:], " ")
			j.createEntry(entry)
		default:
			// the user has provided more than one word. Time to add an entry.
			entry = strings.Join(os.Args[1:], " ")
			j.createEntry(entry)
		}
	}

	// save journal
	j.save()
}

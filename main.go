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

	// loop thought all arguments
	switch args := len(os.Args); {
	case args == 1:
		// user has provided no args. Prompt him to say something
		fmt.Println("Write your entry or provide some args. Available args:")
		fmt.Println(strings.Join(j.Args, "\n"))
		return

	case args > 1:
		// detect the passed arg
		switch strings.ToLower(os.Args[1]) {
		case "--remove":
			fmt.Println("REMOVE")
		case "--view":
			j.showDay(os.Args[2])
		case "--search":
			fmt.Println("SEARCH")
		case "--add":
			// the user has provided the "--add" flag. Time to add an entry.
			entry := strings.Join(os.Args[2:], " ")
			j.createEntry(entry)
			// save journal
			j.save()
		default:
			// the user has provided more than one word. Time to add an entry.
			entry := strings.Join(os.Args[1:], " ")
			j.createEntry(entry)
			// save journal
			j.save()
		}
	}
}

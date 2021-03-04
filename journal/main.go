// made by Lorenzo Rossi - www.lorenzoros.si
// GitHub repo: github.com/lorossi/journal

package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func main() {
	// add flags
	version := flag.Bool("version", false, "show current version")
	use := flag.String("use", "", "use a journal that's not the default one")
	add := flag.String("add", "", "add an entry to the journal. Date format: today, yesterday, weekday (monday-sunday) YYYY-MM-DD, YYYY-MM-DD. You can also set a time in format hh.mm")
	remove := flag.String("remove", "", "remove an entry from the journal. Date format: YYYY-MM-DD or YYYY-MM or YYYY")
	show := flag.String("show", "", "show entries from the journal. Use all to see all. Date format: YYYY-MM-DD or YYYY-MM or YYYY")
	from := flag.String("from", "", "starting date. Only valied if passed with --show, --search or --remove flags and \"all\" argument. Format: YYYY-MM-DD")
	to := flag.String("to", "", "ending date. Only valied if passed with --show, --search or --remove flag and \"all\" argument. Format: YYYY-MM-DD")
	searchkeywords := flag.String("search", "", "search entries by text (both in title and content)")
	searchtags := flag.String("searchtags", "", "search entries by tags")
	searchfields := flag.String("searchfields", "", "search entries by fields")
	printPlaintext := flag.Bool("plaintext", false, "show as plaintext")
	printJSON := flag.Bool("json", false, "show as json")
	tags := flag.Bool("tags", false, "show all used tags")
	fields := flag.Bool("fields", false, "show all used fields")
	encrypt := flag.Bool("encrypt", false, "encrypt journal using AES")
	decrypt := flag.Bool("decrypt", false, "decrypt using AES")
	removePassword := flag.Bool("removepassword", false, "permanently decrypt a journal. This cannot be reversed.")

	flag.Parse()

	// no commands were provided and no text was written
	if flag.NFlag() == 0 && flag.NArg() == 0 {
		flag.PrintDefaults()
		// now exit
		return
	}

	// create empty Journal
	j, e := NewJournal()
	if e != nil {
		printError(e, 2)
		return
	}

	if *version {
		color.Set(color.FgHiGreen)
		fmt.Print("\n\tJournal Version ")
		color.Set(color.FgHiBlue)
		fmt.Print(j.Version, "\n")
		color.Set(color.FgHiGreen)
		fmt.Print("\tGitHub repo: ")
		color.Set(color.FgHiBlue)
		fmt.Print(j.Repo, "\n")

		currentVersion, e := j.getCurrentVersion()

		if e == nil {
			if j.Version != currentVersion {
				color.Set(color.FgHiRed)
				fmt.Print("\tNew version available: ")
				fmt.Print(currentVersion, "\n\n")
			} else {
				color.Set(color.FgHiGreen)
				fmt.Print("\tYou are running the most recent version\n\n")
			}
		} else {
			fmt.Print("\n")
		}

		color.Unset()
		return
	}

	// journal is not using the default filename
	if *use != "" {
		// check if the string ends in .json
		// if not, append it
		if string((*use)[len(*use)-5:]) != ".json" {
			*use += ".json"
		}
		j.setFilename(*use)
	}

	// load from database
	if *decrypt {
		password, e := getPassword("Decryption password:")
		j.setPassword(password)
		if e != nil {
			printError(e, 2)
			return
		}

		e = j.decrypt()
		if e != nil {
			printError(e, 3)
			return
		}
	} else {
		e := j.load()
		if e != nil {
			printError(e, 3)
			return
		}
	}

	// no commands were provided but some text was recognized
	if flag.NFlag() == 0 && flag.NArg() > 0 {
		entry := strings.Join(flag.Args(), " ")
		j.createEntry(entry)
	} else if *add != "" {
		// get text provided by the flag
		// get remainder text
		// concantenate them
		entry := string(*add) + " " + strings.Join(flag.Args(), " ")
		j.createEntry(entry)
	} else if *remove != "" {
		var e error
		if *remove == "all" {
			if *from != "" && *to != "" {
				e = j.removeEntriesBetween(*from, *to)
			} else if *from == "" && *to == "" {
				j.removeAllEntries()
			} else {
				printError(errors.New("wrong parameter with remove flag"), 2)
			}
		} else {
			e = j.removeEntry(*remove)
		}
		if e != nil {
			printError(e, 2)
		}
	} else if *show != "" {
		// get entry by date
		if strings.ToLower(*show) == "all" {
			// check if parameter is "all"
			entries, e := j.getAllEntries()
			if e != nil {
				printError(e, 1)
			} else {
				printEntries(entries, *printPlaintext, *printJSON)
			}
		} else if *from != "" && *to != "" {
			// get entries between dates
			entries, e := j.getEntriesBetween(*from, *to)
			if e != nil {
				printError(e, 1)
			} else {
				printEntries(entries, *printPlaintext, *printJSON)
			}
		} else {
			// check if the parameter is some kind of date
			entries, e := j.showEntries(*show)
			if e != nil {
				printError(e, 1)
			} else {
				printEntries(entries, *printPlaintext, *printJSON)
			}
		}
	} else if *searchkeywords != "" {
		var keywords []string
		// concantenate all the keywords
		keywords = append(keywords, *searchkeywords)
		keywords = append(keywords, flag.Args()...)
		entries, e := j.searchKeywords(keywords)
		if e != nil {
			printError(e, 1)
		} else {
			printEntries(entries, *printPlaintext, *printJSON)
		}
	} else if *searchtags != "" {
		var tags []string
		// concantenate all the keywords
		tags = append(tags, *searchtags)
		tags = append(tags, flag.Args()...)
		entries, e := j.searchTags(tags)
		if e != nil {
			printError(e, 1)
		} else {
			printEntries(entries, *printPlaintext, *printJSON)
		}
	} else if *searchfields != "" {
		var keys []string
		// concantenate all the fields keys
		keys = append(keys, *searchfields)
		keys = append(keys, flag.Args()...)
		entries, e := j.searchFields(keys)
		if e != nil {
			printError(e, 1)
		} else {
			printEntries(entries, *printPlaintext, *printJSON)
		}
	} else if *tags {
		var tags map[string]int
		tags, e := j.getAllTags()
		if e != nil {
			printError(e, 1)
		} else {
			printTags(tags)
		}
	} else if *fields {
		var fields []map[string]string
		fields, e := j.getAllFields()
		if e != nil {
			printError(e, 1)
		} else {
			printFields(fields)
		}
	} else if !(*encrypt || *decrypt || *removePassword) {
		// not a single valid option has been called
		flag.PrintDefaults()
		// now exit
		return
	}

	if *encrypt {
		var password string
		password, e = getPassword("Encryption password:")
		if confirmPassword, _ := getPassword("Confirm password:"); confirmPassword != password {
			j.save()
			printError(errors.New("the two passwords don't match. Saving in plaintext,"), 2)
		} else {
			j.setPassword(password)
			j.encrypt()
		}
	} else if *removePassword {
		e = j.save()
	} else if *decrypt {
		j.encrypt()
	} else {
		e = j.save()
		if e != nil {
			printError(e, 2)
		}
	}

	if e != nil {
		printError(e, 2)
		return
	}
}

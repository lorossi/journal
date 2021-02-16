package main

import (
	"flag"
	"strings"
)

func main() {
	// add flags
	add := flag.String("add", "", "add an entry to the diary. Date format: today:, yesterday:, YYYY-MM-DD")
	remove := flag.String("remove", "", "remove an entry from the diary. Date format: YYYY-MM-DD or YYYY-MM or YYYY")
	view := flag.String("view", "", "view an entry or all entries from the diary. Use all to see all. Date format: YYYY-MM-DD or YYYY-MM or YYYY")
	searchkeywords := flag.String("searchkeywords", "", "search entries by keyword")
	searchtags := flag.String("searchtags", "", "search entries by tags")
	searchfields := flag.String("searchfields", "", "search entries by fields")
	print_plaintext := flag.Bool("plaintext", false, "show as plaintext")
	print_json := flag.Bool("json", false, "show as json")
	hour := flag.String("time", "", "set a time. Only valied if passed with \"add\" flag. Format: hh.mm (24 hour format)")
	tags := flag.Bool("tags", false, "show all tags")
	fields := flag.Bool("fields", false, "show all fields")
	from := flag.String("from", "", "starting date. Only valied if passed with --view --remove flags and \"all\" argument. Format: YYYY-MM-DD")
	to := flag.String("to", "", "ending date. Only valied if passed with --view --remove flag and \"all\" argument. Format: YYYY-MM-DD")
	flag.Parse()

	// no commands were provided and no text was written
	if flag.NFlag() == 0 && flag.NArg() == 0 {
		flag.PrintDefaults()
		// now exit
		return
	}

	// vreate empty Journal
	j := crate_journal()
	// load from database
	e := j.load()

	if e != nil {
		print_error(e, 3)
		return
	}

	// no commands were provided but some text was recognized
	if flag.NFlag() == 0 && flag.NArg() > 0 {
		entry := strings.Join(flag.Args(), " ")
		j.createEntry(entry, "")
	} else if *add != "" {
		// get text provided by the flag
		// get remainder text
		// concantenate them
		entry := string(*add) + " " + strings.Join(flag.Args(), " ")
		j.createEntry(entry, *hour)
	} else if *remove != "" {
		var e error
		if *remove == "all" && *from != "" && *to != "" {
			e = j.removeEntriesBetween(*from, *to)
		} else {
			e = j.removeEntry(*remove)
		}
		if e != nil {
			print_error(e, 2)
		}
	} else if *view != "" {
		// get entry by date
		if strings.ToLower(*view) == "all" {
			// check if parameter is "all"
			entries, e := j.getAllEntries()
			if e != nil {
				print_error(e, 1)
			} else {
				for _, entry := range entries {
					print_entry(entry, *print_plaintext, *print_json)
				}
			}
		} else if *from != "" && *to != "" {
			// get entries between dates
			entries, e := j.getEntriesBetween(*from, *to)
			if e != nil {
				print_error(e, 1)
			} else {
				for _, entry := range entries {
					print_entry(entry, *print_plaintext, *print_json)
				}
			}
		} else {
			// check if the parameter is some kind of date
			entries, e := j.viewEntries(*view)
			if e != nil {
				print_error(e, 1)
			} else {
				for _, entry := range entries {
					print_entry(entry, *print_plaintext, *print_json)
				}
			}
		}
	} else if *searchkeywords != "" {
		var keywords []string
		// concantenate all the keywords
		keywords = append(keywords, *searchkeywords)
		keywords = append(keywords, flag.Args()...)
		entries, e := j.searchKeywords(keywords)
		if e != nil {
			print_error(e, 1)
		} else {
			for _, entry := range entries {
				print_entry(entry, *print_plaintext, *print_json)
			}
		}
	} else if *searchtags != "" {
		var tags []string
		// concantenate all the keywords
		tags = append(tags, *searchtags)
		tags = append(tags, flag.Args()...)
		entries, e := j.searchTags(tags)
		if e != nil {
			print_error(e, 1)
		} else {
			for _, entry := range entries {
				print_entry(entry, *print_plaintext, *print_json)
			}
		}
	} else if *searchfields != "" {
		var keys []string
		// concantenate all the fields keys
		keys = append(keys, *searchfields)
		keys = append(keys, flag.Args()...)
		entries, e := j.searchFields(keys)
		if e != nil {
			print_error(e, 1)
		} else {
			for _, entry := range entries {
				print_entry(entry, *print_plaintext, *print_json)
			}
		}
	} else if *tags {
		var tags map[string]int
		tags, e = j.getAllTags()
		if e != nil {
			print_error(e, 1)
		} else {
			print_tags(tags)
		}
	} else if *fields {
		var fields []map[string]string
		fields, e := j.getAllFields()
		if e != nil {
			print_error(e, 1)
		} else {
			print_fields(fields)
		}
	} else {
		// not a single valid option has been called
		flag.PrintDefaults()
		// now exit
		return
	}

	e = j.save()

	if e != nil {
		print_error(e, 3)
	}
}

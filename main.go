package main

import (
	"flag"
	"strings"
)

func main() {
	// add flags
	add := flag.String("add", "", "add an entry to the diary")
	remove := flag.String("remove", "", "remove an entry from the diary")
	view := flag.String("view", "", "view an entry or all entries from the diary. Use all to see all. Date format: YYYY-MM-DD")
	searchkeywords := flag.String("searchkeywords", "", "search entries by keyword")
	searchtags := flag.String("searchtags", "", "search entries by tags")
	searchfields := flag.String("searchfields", "", "search entries by fields")
	loadmonth := flag.String("loadmonth", "", "load all entries from one month. Date format YYYY-MM")
	loadyear := flag.String("loadyear", "", "load all entries from one year. Date format YYYY")
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
	j.load()

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
		e := j.removeEntry(*remove)
		if e != nil {
			print_error(e)
		}
	} else if *view != "" {
		if strings.ToLower(*view) == "all" {
			entries, e := j.getAllEntries()
			if e != nil {
				print_error(e)
			} else {
				for _, entry := range entries {
					print_entry(entry)
				}
			}
		} else {
			entry, e := j.getEntry(*view)
			if e != nil {
				print_error(e)
			} else {
				print_entry(entry)
			}
		}
	} else if *searchkeywords != "" {
		var keywords []string
		// concantenate all the keywords
		keywords = append(keywords, *searchkeywords)
		keywords = append(keywords, flag.Args()...)
		entries, e := j.searchKeywords(keywords)
		if e != nil {
			print_error(e)
		} else {
			for _, entry := range entries {
				print_entry(entry)
			}
		}
	} else if *searchtags != "" {
		var tags []string
		// concantenate all the keywords
		tags = append(tags, *searchtags)
		tags = append(tags, flag.Args()...)
		entries, e := j.searchTags(tags)
		if e != nil {
			print_error(e)
		} else {
			for _, entry := range entries {
				print_entry(entry)
			}
		}
	} else if *searchfields != "" {
		var keys []string
		// concantenate all the fields keys
		keys = append(keys, *searchfields)
		keys = append(keys, flag.Args()...)
		entries, e := j.searchFields(keys)
		if e != nil {
			print_error(e)
		} else {
			for _, entry := range entries {
				print_entry(entry)
			}
		}
	} else if *loadmonth != "" {
		entries, e := j.getMonth(*loadmonth)
		if e != nil {
			print_error(e)
		} else {
			for _, entry := range entries {
				print_entry(entry)
			}
		}
	} else if *loadyear != "" {
		entries, e := j.getYear(*loadyear)
		if e != nil {
			print_error(e)
		} else {
			for _, entry := range entries {
				print_entry(entry)
			}
		}
	}

	j.save()
}

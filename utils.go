package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

func find_delimiter(entry string, delimiters []string) string {
	for _, e := range entry {
		for _, d := range delimiters {
			if d == string(e) {
				return string(d)
			}
		}
	}

	return ""
}

func remove_multiple_spaces(entry string) string {
	for strings.Contains(entry, "  ") {
		entry = strings.ReplaceAll(entry, "  ", " ")
	}
	return entry
}

func parse_entry(entry, time_format string) (string, time.Time) {
	parsed_day, level := parse_day(entry)
	if level == 0 {
		words := strings.Split(entry, " ")
		return strings.Join(words[1:], " "), parsed_day
	} else {
		return entry, time.Now()
	}
}

func parse_day(entry string) (time_obj time.Time, level int) {
	switch first_word := strings.Split(entry, " ")[0]; strings.ToLower(first_word) {
	case "yesterday":
		// the first word was yesterday. Return today's date MINUS one day
		return time.Now().AddDate(0, 0, -1), 0
	case "today":
		// the first word was today. Return today's date
		return time.Now(), 0
	default:
		// the first word wasn't either yesterday or today.
		// try to parse the date. If it work, remove the first word.
		// If it doesn't work, the date is today (the first word
		// does not indicate the date)
		time_templates := [...]string{"2006-01-02", "2006-01", "2006"}

		for i, template := range time_templates {
			time_obj, e := time.Parse(template, first_word)
			if e == nil {
				return time_obj, i
			}
		}

		// now try matching against weekday
		_, e := time.Parse("Monday", first_word)
		if e == nil {
			now := time.Now()
			for !strings.EqualFold(first_word, now.Weekday().String()) {
				now = now.AddDate(0, 0, -1)
			}
			return now, 0
		}

		return time.Time{}, -1
	}
}

func same_day(date_1, date_2 time.Time) bool {
	return date_1.Format("20060102") == date_2.Format("20060102")
}

func same_month(date_1, date_2 time.Time) bool {
	return date_1.Month() == date_2.Month()
}

func same_year(date_1, date_2 time.Time) bool {
	return date_1.Year() == date_2.Year()
}

func date_between(current, start, end time.Time) bool {
	return current.After(start) && current.Before(end)
}

func print_entry(entry Entry, print_plaintext bool, print_json bool) {
	if print_plaintext {
		// print date
		fmt.Print("[", entry.Timestamp, "] ")
		// print title
		fmt.Print(entry.Title, " ")
		// print content
		fmt.Print(entry.Content, " ")
		// print tags
		if len(entry.Tags) > 0 {
			fmt.Print("+" + strings.Join(entry.Tags, " +"))
		}
		fmt.Print(" ")
		// print fields
		for k, v := range entry.Fields {
			fmt.Print(k, "=", v, " ")
		}
		fmt.Print(" ")
		// end line
		fmt.Println()
	} else if print_json {
		JSON_bytes, _ := json.MarshalIndent(entry, "", "  ")
		fmt.Println(string(JSON_bytes))
	} else {
		// print timestamp
		fmt.Println()
		color.Set(color.FgHiBlue)
		fmt.Print("Date: ")
		color.Unset()
		fmt.Println(entry.Timestamp)

		// print title
		color.Set(color.FgHiGreen)
		fmt.Print("Title: ")
		color.Unset()
		fmt.Println(entry.Title)

		// print content
		color.Set(color.FgHiGreen)
		fmt.Print("Content: ")
		color.Unset()
		fmt.Println(entry.Content)

		// print tags
		color.Set(color.FgHiMagenta)
		fmt.Print("Tags: ")
		color.Unset()
		if len(entry.Tags) > 0 {
			fmt.Print("+" + strings.Join(entry.Tags, " +"))
		}
		fmt.Println()

		// print fields
		color.Set(color.FgHiMagenta)
		fmt.Print("Fields: ")
		color.Unset()
		for k, v := range entry.Fields {
			fmt.Print(k, "=", v, " ")
		}

		// add some spacing
		fmt.Println()
		fmt.Println()
	}
}

func print_tags(tags map[string]int) {
	for k, v := range tags {
		// print key
		color.Set(color.FgHiMagenta)
		fmt.Print(k, " ")
		// print value
		color.Unset()
		fmt.Print(v)
		// end line
		fmt.Println()
	}
}

func print_fields(fields []map[string]string) {
	for _, f := range fields {
		for k, v := range f {
			// print key
			color.Set(color.FgHiMagenta)
			fmt.Print(k, " ")
			// print value
			color.Unset()
			fmt.Print(v)
			// end line
			fmt.Println()
		}
	}
}

func print_error(e error, level int8) {
	switch level {
	case 0:
		color.Set(color.FgHiGreen)
	case 1:
		color.Set(color.FgHiYellow)
	case 2:
		color.Set(color.FgHiRed)
	case 3:
		color.Set(color.BgHiRed)
		color.Set(color.FgHiWhite)
	}
	fmt.Println(e)
}

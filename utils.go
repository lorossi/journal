package main

import (
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

func parse_day_entry(entry, time_format string) (string, time.Time) {
	switch words := strings.Split(entry, " "); strings.ToLower(words[0]) {
	case "yesterday:":
		// the first word was yesterday. Return today's date MINUS one day
		return strings.Join(words[1:], " "), time.Now().AddDate(0, 0, -1)
	case "today:":
		// the first word was today. Return today's date
		return strings.Join(words[1:], " "), time.Now()
	default:
		// the first word wasn't either yesterday or today.
		// try to parse the date. If it work, remove the first word.
		// If it doesn't work, the date is today (the first word
		// does not indicate the date)
		time_obj, e := time.Parse(time_format, words[0])
		if e == nil {
			return strings.Join(words[1:], " "), time_obj
		} else {
			return entry, time.Now()
		}
	}
}

func parse_day(entry, time_format string) time.Time {
	switch first_word := strings.Split(entry, " ")[0]; strings.ToLower(first_word) {
	case "yesterday":
		// the first word was yesterday. Return today's date MINUS one day
		return time.Now().AddDate(0, 0, -1)
	case "today":
		// the first word was today. Return today's date
		return time.Now()
	default:
		// the first word wasn't either yesterday or today.
		// try to parse the date. If it work, remove the first word.
		// If it doesn't work, the date is today (the first word
		// does not indicate the date)
		time_obj, e := time.Parse(time_format, first_word)
		if e == nil {
			return time_obj
		} else {
			// returns zero (epoch time)
			return time.Time{}
		}
	}
}

func same_date(date_1, date_2 time.Time) bool {
	return date_1.Format("20060102") == date_2.Format("20060102")
}

func same_month(date_1, date_2 time.Time) bool {
	return date_1.Month() == date_2.Month()
}

func same_year(date_1, date_2 time.Time) bool {
	return date_1.Year() == date_2.Year()
}

func print_entry(entry Entry) {
	fmt.Println()
	color.Set(color.FgGreen)
	fmt.Print("Date: ")
	color.Unset() // Don't forget to unset
	fmt.Println(entry.Timestamp)

	color.Set(color.FgGreen)
	fmt.Print("Title: ")
	color.Unset()
	fmt.Println(entry.Title)

	color.Set(color.FgGreen)
	fmt.Print("Content: ")
	color.Unset()
	fmt.Println(entry.Content)

	color.Set(color.FgMagenta)
	fmt.Print("Tags: ")
	color.Unset()
	if len(entry.Tags) > 0 {
		fmt.Print("+" + strings.Join(entry.Tags, " +"))
	}
	fmt.Println()

	color.Set(color.FgMagenta)
	fmt.Print("Fields: ")
	color.Unset()
	for k, v := range entry.Fields {
		fmt.Print(k, "=", v, " ")
	}

	fmt.Println()
	fmt.Println()
}

func print_error(e error) {
	color.Set(color.FgYellow)
	fmt.Println(e)
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"golang.org/x/term"
)

func get_password() (password string, e error) {
	fmt.Print("Password: ")
	bytepw, e := term.ReadPassword(int(os.Stdin.Fd()))
	if e != nil {
		return "", errors.New("cannot load password")
	}

	// newline (doesn't get added automatically)
	fmt.Println()
	// pad password to length
	for len(bytepw) < 32 {
		bytepw = append(bytepw, '0')
	}
	// if the password it's too long, chop it
	if len(bytepw) > 32 {
		bytepw = bytepw[:32]
	}

	return string(bytepw), e
}

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
		return strings.Join(words[2:], " "), parsed_day
	} else if level == 1 {
		words := strings.Split(entry, " ")
		return strings.Join(words[1:], " "), parsed_day
	} else {
		return entry, time.Now()
	}
}

func parse_day(entry string) (time_obj time.Time, level int) {
	first_word := strings.Split(entry, " ")[0]
	second_word := strings.Split(entry, " ")[1]

	hour_obj, hour_err := func(hour_str string) (time.Time, error) {
		time_obj, e := time.Parse("15.04", hour_str)
		if e == nil {
			return time_obj, e
		}
		return time.Time{}, e
	}(second_word)

	date_obj, level := func(date_str string) (time.Time, int) {
		switch first_word := strings.Split(entry, " ")[0]; strings.ToLower(first_word) {
		case "yesterday":
			// the first word was yesterday. Return today's date MINUS one day
			return time.Now().AddDate(0, 0, -1), 1
		case "today":
			// the first word was today. Return today's date
			return time.Now(), 1
		default:
			// check the second word, it might be time

			// the first word wasn't either yesterday or today.
			// try to parse the date. If it work, remove the first word.
			// If it doesn't work, the date is today (the first word
			// does not indicate the date)
			time_templates := [...]string{"2006-01-02", "2006-01-02", "2006-01", "2006"}

			for level, template := range time_templates {
				time_obj, e := time.Parse(template, first_word)
				if e == nil {
					return time_obj, level + 1
				}
			}

			// now try matching against weekday
			_, e := time.Parse("Monday", first_word)
			if e == nil {
				now := time.Now()
				for !strings.EqualFold(first_word, now.Weekday().String()) {
					now = now.AddDate(0, 0, -1)
				}
				return now, 1
			}

			return time.Time{}, -1
		}
	}(first_word)

	if hour_err == nil && level != -1 {
		new_date := time.Date(date_obj.Year(), date_obj.Month(), date_obj.Day(), hour_obj.Hour(), hour_obj.Minute(), 0, 0, date_obj.Location())
		return new_date, level - 1
	}
	return date_obj, level
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

func print_entries(entries []Entry, print_plaintext bool, print_json bool) {
	if print_plaintext {
		for _, entry := range entries {
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
		}
	} else if print_json {
		JSON_bytes, _ := json.MarshalIndent(entries, "", "  ")
		fmt.Println(string(JSON_bytes))
	} else {
		for _, entry := range entries {
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
		}
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
	for _, field := range fields {
		for k, v := range field {
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

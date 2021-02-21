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

func getPassword(prompt string) (password string, e error) {
	fmt.Print(prompt, " ")
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

func findDelimiter(entry string, delimiters []string) string {
	for _, e := range entry {
		for _, d := range delimiters {
			if d == string(e) {
				return string(d)
			}
		}
	}

	return ""
}

func removeMultipleSpaces(entry string) string {
	for strings.Contains(entry, "  ") {
		entry = strings.ReplaceAll(entry, "  ", " ")
	}
	return entry
}

func parseEntry(entry, timeFormat string) (string, time.Time) {
	parsedDay, level := parseDay(entry)
	if level == 0 {
		words := strings.Split(entry, " ")
		return strings.Join(words[2:], " "), parsedDay
	} else if level == 1 {
		words := strings.Split(entry, " ")
		return strings.Join(words[1:], " "), parsedDay
	} else {
		return entry, time.Now()
	}
}

func parseDay(entry string) (timeObj time.Time, level int) {
	var dateObj, hourObj time.Time
	var hourErr error

	words := strings.Split(entry, " ")

	// the first word should contain the date
	firstWord := strings.Split(entry, " ")[0]

	dateObj, level = func(date_str string) (time.Time, int) {
		switch date_str {
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
			timeTemplates := [...]string{"2006-01-02", "2006-01", "2006"}

			for level, template := range timeTemplates {
				timeObj, e := time.Parse(template, firstWord)
				if e == nil {
					return timeObj, level + 1
				}
			}

			// now try matching against weekday
			_, e := time.Parse("Monday", firstWord)
			if e == nil {
				now := time.Now()
				for !strings.EqualFold(firstWord, now.Weekday().String()) {
					now = now.AddDate(0, 0, -1)
				}
				return now, 1
			}

			return time.Time{}, -1
		}
	}(firstWord)

	// if there's a second word, check if it contains the hour
	if len(words) > 1 {
		secondWord := words[1]

		hourObj, hourErr = func(hour_str string) (time.Time, error) {
			timeObj, e := time.Parse("15.04", hour_str)
			if e == nil {
				return timeObj, e
			}
			return time.Time{}, e
		}(secondWord)

		// if no error has been found, create the new date with the correct hour
		if hourErr == nil && level == 1 {
			newDate := time.Date(dateObj.Year(), dateObj.Month(), dateObj.Day(), hourObj.Hour(), hourObj.Minute(), 0, 0, dateObj.Location())
			return newDate, 0
		}
	}

	return dateObj, level
}

func sameDay(date1, date2 time.Time) bool {
	return date1.Format("20060102") == date2.Format("20060102")
}

func sameMonth(date1, date2 time.Time) bool {
	return date1.Month() == date2.Month()
}

func sameYear(date1, date2 time.Time) bool {
	return date1.Year() == date2.Year()
}

func dateBetween(current, start, end time.Time) bool {
	return current.After(start) && current.Before(end)
}

func printEntries(entries []Entry, printPlaintext bool, printJSON bool) {
	if printPlaintext {
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
	} else if printJSON {
		JSONBytes, _ := json.MarshalIndent(entries, "", "  ")
		fmt.Println(string(JSONBytes))
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

func printTags(tags map[string]int) {
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

func printFields(fields []map[string]string) {
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

func printError(e error, level int8) {
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
	color.Unset()
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lorossi/colorize"
	"golang.org/x/term"
)

// gets password from terminal, hiding the output
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

// finds the first matching delimiter in list
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

// removes multiple spaces from string
func removeMultipleSpaces(entry string) string {
	for strings.Contains(entry, "  ") {
		entry = strings.ReplaceAll(entry, "  ", " ")
	}
	return entry
}

// parse entry (title and date)
func parseEntry(entry string) (string, time.Time) {
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

// loads day from string
// resoulution: level 0 -> minute, level 1 -> day, level 2 -> month, level 3 -> year, level -1 -> undefined
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
					return timeObj, level
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
		if hourErr == nil && level == 0 {
			newDate := time.Date(dateObj.Year(), dateObj.Month(), dateObj.Day(), hourObj.Hour(), hourObj.Minute(), 0, 0, dateObj.Location())
			return newDate, 0
		}
	}

	return dateObj, level
}

// check if two dates are matching down to the minute
func sameMinute(date1, date2 time.Time) bool {
	return date1.Format("20060102-1504") == date2.Format("20060102-1504")
}

// check if two dates are matching down to the day
func sameDay(date1, date2 time.Time) bool {
	return date1.Format("20060102") == date2.Format("20060102")
}

// check if two dates are matching down to the month
func sameMonth(date1, date2 time.Time) bool {
	return date1.Format("200601") == date2.Format("200601")
}

// check if two dates are matching down to the year
func sameYear(date1, date2 time.Time) bool {
	return date1.Format("2006") == date2.Format("2006")
}

// check if a date is between two other dates
func dateBetween(current, start, end time.Time) bool {
	return current.After(start) && current.Before(end)
}

// print enries according to style
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
			fmt.Print(colorize.BrightBlue("Date: "))
			fmt.Println(entry.Timestamp)

			// print title
			colorize.SetStyle(colorize.BrightGreen("Title: "))
			fmt.Println(entry.Title)

			// print content
			colorize.SetStyle(colorize.BrightGreen("Content: "))
			fmt.Println(entry.Content)

			// print tags
			colorize.SetStyle(colorize.BrightMagenta("Tags: "))
			if len(entry.Tags) > 0 {
				fmt.Print("+" + strings.Join(entry.Tags, " +"))
			}
			fmt.Println()

			// print fields
			colorize.SetStyle(colorize.BrightGreen("Fields: "))
			for k, v := range entry.Fields {
				fmt.Print(k, "=", v, " ")
			}

			// add some spacing
			fmt.Println()
		}
		fmt.Println()
	}
}

// print tags (strings starting with + in entry)
func printTags(tags map[string]int) {
	for k, v := range tags {
		// print key
		colorize.SetStyle(colorize.BrightMagenta(k, " "))
		// print value
		fmt.Print(v)
		// end line
		fmt.Println()
	}
}

// print fields (strings starting with @ in entry)
func printFields(fields []map[string]string) {
	for _, field := range fields {
		for k, v := range field {
			// print key
			colorize.SetStyle(colorize.FgBrightMagenta)
			fmt.Print(k, " ")
			// print value
			colorize.ResetStyle()
			fmt.Print(v)
			// end line
			fmt.Println()
		}
	}
}

// print error
// levels: 0 -> 3, from lowest to highest priority
func printError(e error, level int8) {
	switch level {
	case 0:
		colorize.SetStyle(colorize.FgBrightGreen)
	case 1:
		colorize.SetStyle(colorize.FgBrightYellow)
	case 2:
		colorize.SetStyle(colorize.FgBrightRed)
	case 3:
		colorize.SetStyle(colorize.BgBrightRed)
		colorize.SetStyle(colorize.FgBrightWhite)
	}
	fmt.Println(e)
	colorize.ResetStyle()
}

// print current version
func printVersion(version, repo string) {
	colorize.SetStyle(colorize.FgBrightGreen)
	fmt.Print("\n\tJournal Version ")
	colorize.SetStyle(colorize.FgBrightBlue)
	fmt.Print(version, "\n")
	colorize.SetStyle(colorize.FgBrightGreen)
	fmt.Print("\tGitHub repo: ")
	colorize.SetStyle(colorize.FgBrightBlue)
	fmt.Print(repo, "\n")
	colorize.ResetStyle()

	return
}

// print update
func printUpdate(version, newestVersion string) {
	if version != newestVersion {
		colorize.SetStyle(colorize.FgBrightRed)
		fmt.Print("\tNew version available: ")
		fmt.Print(newestVersion, "\n\n")
	} else {
		colorize.SetStyle(colorize.FgBrightGreen)
		fmt.Print("\tYou are running the most recent version\n\n")
	}
}

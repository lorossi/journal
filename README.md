# Journal
*A simple CLI journal made in Go*

## Goals
I want to implement a simple CLI (command line interface) journal in Go, as a way to excercise and learn it.

Wanted features:
  1. Multiple databases
  2. Encryption
  3. Add tags to each Entry
  4. Add pairs of key/value fields to each entry
  5. Show all entries
  6. Filter entries by date, tag, fields, text
  7. MUST BE SIMPLE TO USE!


# Basic usage

## Add entry
Add entry for today:

`journal --add Dear diary, today I was so tired...`

Or skip the parameter and just write the entry:

`journal Dear Diary, today I was so tired...`

You can also specify the day by writing `today`

`journal Today i started using this new journal software! It works really good.`

Add entry for yesterday:

`journal --add yesterday Dear diary, today i studied so much...`

`journal yesterday Dear diary, today i studied so much...`

Add entry for arbitrary date:

`journal --add 2020-02-15 Dear diary, today I read about a strange flu in China. I'm sure it's going to be nothing!`

`journal 2020-02-15 Dear diary, today I read about a strange flu in China. I'm sure it's going to be nothing!`

### Tags
Write *tags* by simply adding a `+` sign before the tag. Example:

`journal Such an exciting day! I went to Disneyland. +fun +happiness`

will store `tag` and `happiness` as tags for today's entry. Of course tags can be used in combination of the previous settings.

### Fields
*Fields* are pairs of key/value. Write them by adding a `@` before the key and `=` before the value. Example:

`journal Today i ran so much! @run=10km @minutes=30`

will store `run: 10km` in todays entry. Of course fields can be used in combination of the previous settings.

### Time
Set a different time (24 hours format) than now for an entry:

`journal --time 06.55 --add today: I just woke up! I totally did not set the time later`

`journal --time 15.30 --add yesterday: i went to bed early!`

`journal --time 9.00 --add 2020/07/03 to the judge: i totally was at home`

## View entry (or multiple entries)
View an entry for an arbitrary date:

`journal --view 2020-02-15`

View all entries from one month or from one year:

`journal --view 2020-01` `journal --view 2020`

View all entries:

`journal --view all`

### View entry between two dates
View entry between two dates (inclusive):

`journal --view all --from 2020-01-01 --to 2021-06-01`

## Remove entry
Remove entry for today:

`journal --remove today`

Remove entry for yesterday:

`journal --remove yesterday`

Remove entry for arbitrary date:

`journal --remove 2020-02-15`

Remove all entries from one month or from one year:

`journal --remove 2020-01` `journal --remove 2020`

### Remove entry between two dates
Remove entry between two dates (inclusive):

`journal --remove all --from 2020-01-01 --to 2021-06-01`

## Search entries by keyword
The keywords will be matched against words in the title and the content of each entry. If an entry matches ANY of the keywords, it will be shown.

Search "skiing":

`journal --searchkeywords skiing`

Search "lake" and "sushi":

`journal --searchkeywords lake sushi`

## Search entries by tag
The tag will be matched against the ones stored in each entry. If an entry matches ANY of the tags, it will be shown.

Search tag "fun":

`journal --searchtags fun`

Search tags "airplane" and "ferry":

`journal --searchtags lake sushi`

### Get all tags
Get all tags and their total usage:

`journal --tags`

## Search entries by field
The field will be matched against the ones stored in each entry. If an entry matches ANY of the fields keys, it will be shown.

Search field with key `pushups`

`journal --searchfields pushups`

Search fields with key `burpess` and `slices_of_cake`:

`journal --searchfields burpees slices_of_cake`

### Get all fields
Get all used fields and their relative values:

`--journal --fields`

# Credits
Thanks to [faith](github.com/fatih) for his [color](github.com/fatih/color) package.

This project is distributed under CC 4.0 License.

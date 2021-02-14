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

You can also specify the day by writing `today:`

`journal Today: i started using this new journal software! It works really good.`

Add entry for yesterday:

`journal --add yesterday: Dear diary, today i studied so much...`

`journal yesterday: Dear diary, today i studied so much...`

Add entry for arbitrary date:

`journal --add 2020-02-15 Dear diary, today I read about a strange flu in China. I'm sure it's going to be nothing!`

`journal 2020-02-15 Dear diary, today I read about a strange flu in China. I'm sure it's going to be nothing!`

### Tags
Write tags by simply adding a `+` sign before the tag. Example:

`journal Such an exciting day! I went to Disneyland. +fun +happiness`

will store `tag` and `happiness` as tags for today's entry. Of course tags can be used in combination of the previous settings.

## View entry (or multiple entries)
View an entry for an arbitrary date:

`journal --view 2020-02-15`

View all entries:

`journal --view all`

## Remove entry
Remove entry for today:

`journal --remove today`

Remove entry for yesterday:

`journal --remove yesterday`

Remove entry for arbitrary date:

`journal --remove 2020-02-15`

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

## Load all entries from one month
`journal --loadmonth 2020-02`

## Load all entries from one year
`journal --loadyear 2020`

# Credits
Thanks to [faith](github.com/fatih) for his [color](github.com/fatih/color) package

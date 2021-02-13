# Journal
~A simple CLI journal made in Go~

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

## Entries
Add entry for today:
`journal --add Dear diary, today I was so tired...`

Or skip the parameter and just write the entry
`journal Dear Diary, today I was so tired...`

Add entry for yesterday:
`journal --add yesterday: Dear diary, today i studied so much...`

`journal yesterday: Dear diary, today i studied so much...`

Add entry for arbitrary date:
`journal --add 2020-02-15 Dear diary, today I read about a strange flu in China. I'm sure it's going to be nothing!`

`journal 2020-02-15 Dear diary, today I read about a strange flu in China. I'm sure it's going to be nothing!`

## Tags
Write tags by simply adding a `+` sign before the tag. Example:

`journal Such an exciting day! I went to Disneyland. +fun +happiness`

will store `tag` and `happiness` as tags for today's entry. Of course tags can be used in combination of the previous settings.

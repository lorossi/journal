# Journal
![](https://img.shields.io/github/repo-size/lorossi/go-journal?style=flat-square)
![](https://img.shields.io/maintenance/yes/2021?style=flat-square)
![](https://img.shields.io/github/last-commit/lorossi/go-journal/main?style=flat-square)
![](https://img.shields.io/github/go-mod/go-version/lorossi/go-journal?style=flat-square)

*A simple CLI journal made in Go*

## Installation
Clone the repo and use the pre compiled binaries inside the `binaries` folder (or download the [latest release](https://github.com/lorossi/go-journal/releases/latest)).

- On *Linux*, use the installer by running `sudo sh installer.sh` to move the binary inside the `PATH` folder in order to run it from everywhere.
- On *Windows*, use the installer by running `installer.bat` from the command prompt/powershell of simply double click on the file
- On *MacOs* may god help you because for sure I can't.

Otherwise, clone the repo and build it from source.

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

`journal today 12.32 I just woke up! I totally did not set the time later`

`journal yesterday 7.24 i went to bed early!`

`journal 2020/07/03 9.00 to the judge: i totally was at home`

## View entry (or multiple entries)
View an entry for an arbitrary date:

`journal --show 2020-02-15`

View all entries from one month or from one year:

`journal --show 2020-01` `journal --show 2020`

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

Remove all entries from the diary:

`journal --remove all`

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

`journal --searchtags airplane ferry`

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

`journal --fields`

## Password protection
The program supports password protection with the AES Encryption algorithm.

While decrypted, the database **will never** be stored in plaintext on your pc.

### Encryption
Encrypt a clear database by using the flag `--encrypt`. You will be asked for a password. **Save it** because it won't be stored and if you lose it there's no way of unlocking your journal again.

`journal --encrypt`

### Decryption
Like in encryption, to decrypt a journal in order to write/read on it, use the flag `--decrypt`. You will be asked for a password.

In case you want to change the password of your journal, use both `--encrypt` and `--decrypt` flags.

### Password removal / change
If you want to remove the password from your journal, you first have to decrypt it by providing a correct password.

`journal --decrypt --removepassword`

In order to change the password, you have to decrypt and then encrypt it again

`journal --decrypt --encrypt`

## Output formatting
The output can be formatted either in JSON or plain text by using the correct flags.

`journal --view all --json`

`journal --view all --plaintext`

## Help
Use the flag `-h` or `--help` to get a list of all the available options.

# Credits
Thanks to [faith](github.com/fatih) for his [color](github.com/fatih/color) package.

This project is distributed under CC 4.0 License.

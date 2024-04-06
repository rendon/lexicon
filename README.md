# Lexicon
A command line lexicon (dictionary).

## Prerequisites
- Go
- SQLite (optional, for data browsing)

## Setup
Create the SQLite database and move the file to its final destination.
```sh
sqlite3 lexicon.sqlite < schema.sql
mv lexicon.sqlite /path/to/db/
```

Set the dictionaryapi.com key in your environment:
```sh
set DICTIONARY_API_KEY="..."
```

Build the binary:
```sh
go build
```

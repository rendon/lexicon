# Status
The build is broken, and re-designing the program.
What's changing:
1. We'll have two places where to store the data
  a. A sqlite database for local usage only
  b. An API for distributed usage (i.e., can be used from multiple computers)
2. I've defined a Dictionary interface with the functions that each dictionary type (sqlite, api, etc.) must implement
3. The source of the definitions will continue to be dictionaryapi.com and it will be queried by this program
1. The WOD will come from a separate source (the wod service), not part of the dictionary API

-- # Lexicon
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

## How to use
Simply execute the binary like so:
```sh
./lexicon
```

You can also install the binary so it's available everywhere. Just make sure the `$GOPATH/bin/` directory is part of your `$PATH`:
```sh
go install
```

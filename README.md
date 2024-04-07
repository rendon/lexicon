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

## How to use
Simply execute the binary like so:
```sh
./lexicon
```

You can also install the binary so it's available everywhere. Just make sure the `$GOPATH/bin/` directory is part of your `$PATH`:
```sh
go install
```

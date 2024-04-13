#!/usr/bin/env bash
if [[ $2 == "--timestamp" ]]; then
	sqlite3 $DATA_SOURCE_NAME "SELECT name, datetime(createdAt / 1000, 'unixepoch') FROM lexicon WHERE createdAt >= strftime('%s', 'now') * 1000 - 86400000"
else
	sqlite3 $DATA_SOURCE_NAME "SELECT name FROM lexicon WHERE createdAt >= strftime('%s', 'now') * 1000 - 86400000"
fi

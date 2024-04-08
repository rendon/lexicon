#!/usr/bin/env bash
sqlite3 $DATA_SOURCE_NAME "SELECT count(*) FROM lexicon"

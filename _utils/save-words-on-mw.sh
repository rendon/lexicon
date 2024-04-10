#!/usr/bin/env bash

FILE=$1

while read word; do
  escaped_word=$(echo $word | sed 's/ /+/g')
  RES=$(curl 'https://www.merriam-webster.com/lapi/v1/wordlist/save' --compressed -X POST\
      -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0'\
      -H 'Accept: application/json, text/javascript, */*; q=0.01'\
      -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8'\
      -H "Cookie: $MERRIAM_WEBSTER_COOKIE"\
      --data-raw 'word='$escaped_word'&type=d' 2>/dev/null)

  if [[ $RES =~ "already exists" ]]
  then
    echo "\"$word\" already saved"
  else
    echo "Added \"$word\""
  fi
done <$FILE

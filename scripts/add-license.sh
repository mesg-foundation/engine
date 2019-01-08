#!/bin/bash -e

copyrighgrep="Copyright 2018 MESG Foundation"
copyright=$(cat LICENSE | awk "/$copyrighgrep/ { matched = 1 } matched" | sed 's/^  //' | sed 's/^/\/\//')

for file in $(find . -name '*.go' -not -path '**/vendor/*');
do
  if ! grep -q "$copyrighgrep" $file
  then
    echo $"$copyright 

$(cat $file)" > $file
  fi
done

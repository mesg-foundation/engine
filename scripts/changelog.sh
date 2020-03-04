#!/bin/bash

[ -z "$1" ] && echo "milestone required" && exit

MILESTONE="$1"

LABELS=("breaking change" "release:add" "release:change" "release:fix" "release:remove" "dependencies")
CATEGORIES=("Breaking Changes" "Added" "Changed" "Fixed" "Removed" "Dependencies")

PR=$(hub pr list -s merged -f "%Mt %L|- ([%i](%U)) %t.%n" --sort-ascending | grep $MILESTONE)

printf "## [$MILESTONE](https://github.com/mesg-foundation/engine/releases/tag/$MILESTONE)\n\n"

for i in "${!LABELS[@]}"; do 
  LABEL=${LABELS[$i]}
  CATEGORY=${CATEGORIES[$i]}
  printf "#### $CATEGORY\n\n"
  echo "$PR" | grep "$LABEL" | cut -d'|' -f2
  printf "\n"
done

#!/usr/bin/env bash

echo "replacing $1 with $2"

for file in $(ag "$1" -l); do
  gsed -i "s|$1|$2|g" "$file"
done

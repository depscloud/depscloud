#!/usr/bin/env bash

readonly version="${VERSION:-"next"}"
readonly sha="${GITHUB_SHA:-"$(git rev-parse HEAD)"}"

# Set metadata in package.json

cat package.json | \
  jq --arg version "$version" '.version = $version' | \
  jq --arg version "$version" '.meta.version = $version' | \
  jq --arg sha "$sha" '.meta.revision = $sha' | \
  jq --arg date "$(date +%Y-%m-%dT%T)" '.meta.date = $date' > package-new.json

mv package-new.json package.json

## Set properly in package-lock.json

cat package-lock.json | \
  jq --arg version "$version" '.version = $version' > package-lock-new.json

mv package-lock-new.json package-lock.json

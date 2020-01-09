#!/usr/bin/env bash

set -e -o pipefail

git checkout gh-pages
mv incubator-dist/* incubator/
git add .
git commit -m "upload latest charts"
git push

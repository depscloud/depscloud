#!/usr/bin/env bash

git add .
git stash
git checkout gh-pages
git stash pop
git commit -m "upload latest charts"
git push

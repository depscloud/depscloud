#!/usr/bin/env bash

git checkout gh-pages
mv incubator-dist/* incubator/
git add .
git commit -m "upload latest charts"
git push

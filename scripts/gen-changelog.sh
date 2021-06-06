#!/usr/bin/env bash

since="$(git tag -l | tail -n 2 | head -n 1)"

cat <<EOF
## Features

$(git diff --name-only "${since}" | grep ^changelog/ | xargs -I{} bash -c "cat {}; echo ''")

## Fixes

$(git log --format="- %s [%h](https://github.com/depscloud/depscloud/commit/%H)" HEAD...${since} | grep 'fix:')

## Contributors

$(git log --format="- [%an](https://github.com/)" HEAD...${since} | grep -v bot | sort | uniq)
EOF

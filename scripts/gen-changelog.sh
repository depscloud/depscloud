#!/usr/bin/env bash

# produces a single markdown blob containing the changelog information
git diff --name-only "$(git tag -l | tail -n 2 | head -n 1)" | \
    grep ^changelog/ | \
    xargs -I{} bash -c "cat {}; echo ''"

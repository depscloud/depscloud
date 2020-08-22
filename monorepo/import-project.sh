#!/usr/bin/env bash

readonly branch=$(git branch | grep '*' | awk '{print $2}')

for repo in $@; do
    if [[ ! -d "$repo" ]]; then
        echo "[subtree] setting up $repo on $branch"
        git subtree add --prefix $repo git@github.com:depscloud/$repo.git $branch
    else
        echo "[subtree] updating $repo on $branch"
        git subtree pull --prefix $repo git@github.com:depscloud/$repo.git $branch
    fi
done

#!/usr/bin/env bash

readonly branch=$(git branch | grep '*' | awk '{print $2}')

function update_subtree() {
    for repo in $@; do
        if [[ -z "$(git remote | grep $repo)" ]]; then
            echo "[remote] setting up $repo"
            git remote add -f $repo git@github.com:depscloud/$repo.git
        else
            echo "[remote] updating $repo on $branch"
            git fetch $repo $branch
        fi

        if [[ ! -d "$repo" ]]; then
            echo "[subtree] setting up $repo on $branch"
            git subtree add --prefix $repo $repo $branch
        else
            echo "[subtree] updating $repo on $branch"
            git subtree pull --prefix $repo $repo $branch
        fi
    done
}

update_subtree \
    api \
    indexer \
    dockerfiles \
    deploy \
    tracker \
    extractor \
    gateway \
    cli

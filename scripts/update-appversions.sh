#!/usr/bin/env bash

set -e -o pipefail

function github_latest() {
    version=$(curl -sSL "https://api.github.com/repos/$1/releases/latest" | yq r - tag_name)
    echo "${version#"v"}"
}

function update_appversion() {
    app=$1

    chart_latest_app_version=$(yq r charts/${app}/Chart.yaml appVersion)
    latest_app_version=$(github_latest depscloud/${app})

    if [[ "${chart_latest_app_version}" != "${latest_app_version}" ]]; then
        echo "Updating appVersion for ${app}"

        yq w -i charts/${app}/Chart.yaml appVersion "${latest_app_version}"
    fi
}

update_appversion extractor
update_appversion gateway
update_appversion indexer
update_appversion tracker

#!/usr/bin/env bash

set -e -o pipefail

function compute_next_version() {
    echo $1 | awk -F. -v OFS=. 'NF==1{print ++$NF}; NF>1{if(length($NF+1)>length($NF))$(NF-1)++; $NF=sprintf("%0*d", length($NF), ($NF+1)%(10^length($NF))); print}'
}

function update_version() {
    app=$1
    
    if [[ ! -z "$(git status -s | grep "charts/${app}/")" ]]; then
        echo "Updating version for ${app}"

        last_version=$(yq r charts/${app}/Chart.yaml version)
        next_version=$(compute_next_version ${last_version})
        
        yq w -i charts/${app}/Chart.yaml version "${next_version}"
        yq w -i charts/depscloud/Chart.yaml "dependencies.(name==${app}).version" "${next_version}"
    fi
}

update_version extractor
update_version gateway
update_version indexer
update_version tracker
update_version depscloud

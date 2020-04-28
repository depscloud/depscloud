#!/usr/bin/env bash

readonly base_url="https://deps-cloud.github.io/charts"

# start off with the current gh-pages
rm -rf dist/
git clone -q --depth 1 -b gh-pages https://github.com/deps-cloud/charts.git dist

# port in any new templates
for chart in $(find . -iname Chart.yaml | xargs dirname | cut -c 3-); do
    echo "Packaging ${chart}"
    destination=dist/$(dirname ${chart})
    mkdir -p ${destination}

    helm dependency update ${chart}
    helm package ${chart} -d ${destination}
done

for repo in "incubator" "stable"; do
    echo "Indexing ${repo}"
    repo_url="${base_url}/${repo}"

    helm repo index dist/${repo} --url "${repo_url}"
done

cp README.md dist/README.md

# generate depscloud simple k8s deployment
mkdir -p dist/deploy/
helm template depscloud ./stable/depscloud/ \
  --set indexer.externalConfig.secretRef.name="depscloud-indexer" \
  --set tracker.externalStorage.secretRef.name="depscloud-tracker" \
  --namespace depscloud-system > dist/deploy/depscloud-system.yaml

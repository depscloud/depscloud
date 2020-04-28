#!/usr/bin/env bash

set -e -o pipefail

readonly base_url="https://deps-cloud.github.io/deploy"
readonly repo_url="https://github.com/deps-cloud/deploy.git"

readonly in=charts
readonly out=public

rm -rf "${out}"
#git clone -q --depth 1 -b gh-pages "${repo_url}" "${out}"

readonly docker_path="${out}/docker/"
readonly k8s_path="${out}/k8s/"
readonly charts_path="${out}/charts"

mkdir -p "${docker_path}"
mkdir -p "${k8s_path}"
mkdir -p "${charts_path}"

## generate helm repository

for chart in "${in}"/*; do
  echo "Linting Helm Chart ${chart}"
  helm lint "${chart}" 1>/dev/null
done

for chart in "${in}"/*; do
  echo "Packaging Helm Chart ${chart}"
  helm dependency update "${chart}" 1>/dev/null
  helm package "${chart}" -d "${charts_path}" 1>/dev/null
done

helm repo index "${charts_path}" --url "${base_url}/charts"

## generate k8s manifests

echo "Packaging Manifests stable/depscloud"
helm template  depscloud ${in}/depscloud/ \
  --set indexer.externalConfig.secretRef.name="depscloud-indexer" \
  --set tracker.externalStorage.secretRef.name="depscloud-tracker" \
  --namespace depscloud-system > "${k8s_path}/depscloud-system.yaml"

## TODO: generate docker-compose

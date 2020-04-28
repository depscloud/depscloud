#!/usr/bin/env bash

set -e -o pipefail

readonly base_url="https://deps-cloud.github.io/deploy"
readonly repo_url="https://github.com/deps-cloud/deploy.git"

readonly in=charts
readonly out=public

rm -rf "${out}"
git clone -q --depth 1 -b gh-pages "${repo_url}" "${out}"

readonly docker_path="${out}/docker"
readonly k8s_path="${out}/k8s"
readonly charts_path="${out}/charts"

rm -rf "${docker_path}"
rm -rf "${k8s_path}"

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

tmp=$(mktemp -d)
trap "rm -rf '${tmp}'" EXIT

echo "Packaging Manifests bitnami/mysql"
helm repo add bitnami https://charts.bitnami.com/bitnami 1>/dev/null
helm template mysql bitnami/mysql \
  --set db.user=user \
  --set db.password=password \
  --set db.name=depscloud \
  --set volumePermissions.enabled=true \
  --namespace depscloud-system > "${tmp}/mysql.yaml"

cat <<EOF > "${tmp}/kustomization.yaml"
namespace: depscloud-system
resources:
- mysql.yaml
EOF

kubectl kustomize "${tmp}" > "${k8s_path}/mysql.yaml"

echo "Packaging Manifests stable/depscloud"
helm template depscloud ${in}/depscloud/ \
  --set indexer.externalConfig.secretRef.name="depscloud-indexer" \
  --set tracker.externalStorage.secretRef.name="depscloud-tracker" \
  --namespace depscloud-system > "${k8s_path}/depscloud-system.yaml"

## copy in README

cp "README.md" "${out}/README.md"

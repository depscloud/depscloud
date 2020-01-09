#!/usr/bin/env bash

readonly os=$(uname | tr '[:upper:]' '[:lower:]')

trap "rm -rf helm.tar.gz ${os}-amd64" EXIT

curl -sSL -o helm.tar.gz https://get.helm.sh/helm-v3.0.2-${os}-amd64.tar.gz
tar zxf helm.tar.gz

git checkout gh-pages

helm repo index incubator/

git add incubator/index.yaml
git commit -m "index latest charts"
git push

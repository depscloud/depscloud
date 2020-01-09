#!/usr/bin/env bash

readonly os=$(uname | tr '[:upper:]' '[:lower:]')

trap "rm -rf helm.tar.gz ${os}-amd64" EXIT

curl -o helm.tar.gz https://get.helm.sh/helm-v3.0.2-${os}-amd64.tar.gz
tar zxvf helm.tar.gz

for chart in $(ls -1 incubator/); do
  ./${os}-amd64/helm package incubator/${chart} -d incubator
done

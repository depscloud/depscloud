pushd ./packages/depscloud-api-nodejs
npm version ${1}
cd ../..

export VERSION="v$(jq -r .version ./packages/depscloud-api-nodejs/package.json)"
echo "version: ${VERSION}"
git commit -a -m "${VERSION}"
git tag -a "${VERSION}" -m "${VERSION}"

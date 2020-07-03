set -e -o pipefail

readonly ROOT_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )" )"
readonly NODEJS_PACKAGE_DIR="${ROOT_DIR}/packages/depscloud-api-nodejs"

cp -R ${ROOT_DIR}/proto/src/depscloud_api/* ${NODEJS_PACKAGE_DIR}
pushd ${NODEJS_PACKAGE_DIR}
[[ -d node_modules/ ]] || npm install
npm run generate
popd

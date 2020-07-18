set -e -o pipefail

if [[ $# -lt 1 ]]; then
  echo "expected at least one argument"
  exit 1
fi

readonly ROOT_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )" )"
readonly NODEJS_PACKAGE_DIR="${ROOT_DIR}/packages/depscloud-api-nodejs"

protoc \
  -I=${ROOT_DIR}/proto/src \
  --dts_out=${NODEJS_PACKAGE_DIR} \
 ${1}

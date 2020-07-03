set -e -o pipefail

readonly ROOT_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )" )"
readonly SWAGGER_DIR="${ROOT_DIR}/swagger"

pushd ${SWAGGER_DIR}
go-bindata -fs -pkg swagger -o swagger.go $(find . -iname *.swagger.json)
popd

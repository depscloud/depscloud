set -e -o pipefail

if [[ $# -lt 1 ]]; then
  echo "expected at least one argument"
  exit 1
fi

dir=src
if [[ $# -gt 1 ]]; then
  dir=${2}
fi

readonly ROOT_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )" )"
readonly PROTO_DIR="${ROOT_DIR}/proto/${dir}"

pushd "${PROTO_DIR}"
for file in $(find . -name *.proto | cut -c 3-); do
  echo "compiling ${file} using ${1}"
  bash ${ROOT_DIR}/scripts/${1}.sh "${file}"
done
popd

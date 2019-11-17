#!/usr/bin/env bash

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [[ "${ARCH}" == "x86_64" ]]; then
    ARCH=amd64
elif [[ "${ARCH}" == "i386" ]]; then
    ARCH=386
elif [[ "${ARCH}" == "aarch64" ]]; then
    ARCH=arm64
else
    ARCH=arm
fi

if [[ $# -lt 1 ]]; then
cat <<EOF
  usage: download-health-probe <version>
EOF
exit 1
fi

readonly USR="grpc-ecosystem"
readonly APP="grpc-health-probe"
readonly APP_VERSION="${1}"

cat <<EOF
user:     ${USR}
app:      ${APP}
version:  ${APP_VERSION}
os:       ${OS}
arch:     ${ARCH}
EOF

curl -L -O "https://github.com/${USR}/${APP}/releases/download/v${APP_VERSION}/grpc_health_probe-${OS}-${ARCH}"
chmod 755 grpc_health_probe-${OS}-${ARCH}
mv grpc_health_probe-${OS}-${ARCH} /usr/bin/grpc_health_probe

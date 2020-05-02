#!/usr/bin/env bash
mkdir -p certs

function gencert() {
    process=${1}

    openssl req -new -newkey rsa:4096 \
        -keyout certs/${process}.key -out certs/${process}.csr \
        -nodes -subj "/CN=${process}"

    openssl x509 -req -sha256 -days 365 \
        -in certs/${process}.csr \
        -CA certs/ca.crt -CAkey certs/ca.key \
        -set_serial 01 -out certs/${process}.crt
}

openssl req -x509 -sha256 -newkey rsa:4096 \
    -keyout certs/ca.key -out certs/ca.crt \
    -days 356 -nodes -subj '/CN=depscloud-ca'

gencert "extractor"
gencert "tracker"
gencert "gateway"
gencert "indexer"
gencert "localhost"
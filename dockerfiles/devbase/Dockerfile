ARG GO_VERSION=1.16
ARG ALPINE_VERSION=3.12

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION}

ENV GO111MODULE on

RUN apk -U upgrade && apk add build-base git ca-certificates sqlite bash curl

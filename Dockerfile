FROM alpine:3.9

RUN apk update && apk add curl
COPY ./install-depscloud-binary /usr/bin/install-depscloud-binary

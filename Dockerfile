FROM debian:stretch-slim

RUN apt-get update -y && apt-get install -y curl
COPY ./install-depscloud-binary /usr/bin/install-depscloud-binary

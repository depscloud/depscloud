ARG BINARY=""

FROM ocr.sh/depscloud/devbase:latest AS builder

WORKDIR /go/src/depscloud

COPY go.mod .
COPY go.sum .
COPY Makefile .
RUN make build-deps deps

ARG VERSION=""
ARG GIT_SHA=""
ARG BINARY

COPY internal/ internal/
COPY services/${BINARY}/ services/${BINARY}/
RUN make install-${BINARY} VERSION="${VERSION}" GIT_SHA="${GIT_SHA}"

FROM ocr.sh/depscloud/base:latest
ARG BINARY

COPY --from=builder /go/bin/${BINARY} /usr/bin/${BINARY}
RUN ln -s /usr/bin/${BINARY} /usr/bin/entrypoint

WORKDIR /home/depscloud
USER 13490:13490

ENTRYPOINT [ "/usr/bin/entrypoint" ]

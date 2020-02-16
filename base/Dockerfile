ARG RUNTIME_USER=deps-cloud
ARG RUNTIME_UID=13490

FROM alpine:3.11

RUN apk update && apk add -u --no-cache ca-certificates tzdata

COPY etc_passwd /etc/passwd

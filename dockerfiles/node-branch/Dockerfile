ARG NODE_VERSION=16

FROM node:${NODE_VERSION}
RUN apt-get update && apt-get install -y jq

WORKDIR /home/depscloud

ARG VERSION=""
ARG GIT_SHA=""
ARG BINARY

COPY services/${BINARY}/package.json .
COPY services/${BINARY}/package-lock.json .

RUN npm install

COPY services/${BINARY}/ .

RUN npm run build && \
    npm run prepackage

USER 13490:13490

ENTRYPOINT [ "npm", "start", "--" ]

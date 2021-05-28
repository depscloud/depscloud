ARG NODE_VERSION=16

FROM node:${NODE_VERSION}

WORKDIR /home/depscloud

ARG VERSION=0.2.10
ARG BINARY=""

ADD https://github.com/depscloud/depscloud/releases/download/v${VERSION}/${BINARY}-${VERSION}.tar.gz ${BINARY}.tar.gz
RUN tar -zxvf ${BINARY}.tar.gz && rm ${BINARY}.tar.gz

RUN npm install --production

USER 13490:13490

ENTRYPOINT [ "npm", "start", "--" ]

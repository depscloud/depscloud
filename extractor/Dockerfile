ARG NODE_VERSION=12

FROM node:${NODE_VERSION}

WORKDIR /home/depscloud

ARG VERSION=0.2.10
ADD https://github.com/depscloud/depscloud/releases/download/v${VERSION}/extractor-${VERSION}.tar.gz extractor.tar.gz
RUN tar -zxvf extractor.tar.gz && rm extractor.tar.gz

RUN npm install --production

USER 13490:13490

ENTRYPOINT [ "npm", "start", "--" ]

FROM node:10

ARG VERSION=0.1.2
ARG HEALTH_PROBE_VERSION=0.3.1

RUN curl -L -o extractor.zip https://github.com/deps-cloud/extractor/archive/v${VERSION}.zip && \
    unzip extractor.zip && \
    mv extractor-${VERSION} /app

COPY download-health-probe.sh /usr/bin/download-health-probe

RUN download-health-probe ${HEALTH_PROBE_VERSION} && \
    rm -rf /usr/bin/download-health-probe

WORKDIR /app

RUN npm install && npm run build

ENTRYPOINT [ "npm", "start", "--" ]

FROM node:10

COPY . /app

WORKDIR /app

RUN rm -rf node_modules && npm install
RUN npm run lint && npm run test && npm run build

ENTRYPOINT [ "npm", "start" ]

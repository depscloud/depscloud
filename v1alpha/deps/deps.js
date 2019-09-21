const path = require('path');
const protoLoader = require('@grpc/proto-loader');
const grpc = require('grpc');

const packageDefinition = protoLoader.loadSync(
    path.join(__dirname, 'deps.proto'),
    {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true,
        includeDirs: [
            __dirname
        ]
    }
);

const descriptor = grpc.loadPackageDefinition(packageDefinition);

module.exports = descriptor.cloud.deps.extractor.api;

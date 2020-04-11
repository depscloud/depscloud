const path = require('path');
const protoLoader = require('@grpc/proto-loader');
const grpc = require('grpc');

const packageDefinition = protoLoader.loadSync(
    [
        path.join(__dirname, "v1alpha", "deps", "deps.proto"),
        path.join(__dirname, "v1alpha", "extractor", "extractor.proto"),
        path.join(__dirname, "v1alpha", "schema", "schema.proto"),
        path.join(__dirname, "v1alpha", "store", "store.proto"),
        path.join(__dirname, "v1alpha", "tracker", "tracker.proto"),
    ],
    {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true,
        includeDirs: [
            __dirname,
        ]
    }
);

const descriptor = grpc.loadPackageDefinition(packageDefinition);

module.exports = descriptor.cloud.deps.api;

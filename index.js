const path = require('path');
const protoLoader = require('@grpc/proto-loader');
const grpc = require('grpc');

const filenames = [];

filenames.push(path.join(__dirname, "v1alpha", "tracker", "tracker.proto"));
filenames.push(path.join(__dirname, "v1alpha", "extractor", "extractor.proto"));
filenames.push(path.join(__dirname, "v1alpha", "schema", "schema.proto"));
filenames.push(path.join(__dirname, "v1alpha", "deps", "deps.proto"));
filenames.push(path.join(__dirname, "v1alpha", "store", "store.proto"));

const packageDefinition = protoLoader.loadSync(
    filenames,
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

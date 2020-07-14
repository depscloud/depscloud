const path = require('path');
const protoLoader = require('@grpc/proto-loader');
const grpc = require('@grpc/grpc-js');

const filenames = [];

filenames.push(path.join(__dirname, "node_modules", "protobufjs", "google", "api", "http.proto"));
filenames.push(path.join(__dirname, "node_modules", "protobufjs", "google", "api", "annotations.proto"));
filenames.push(path.join(__dirname, "node_modules", "protobufjs", "google", "protobuf", "api.proto"));
filenames.push(path.join(__dirname, "node_modules", "protobufjs", "google", "protobuf", "source_context.proto"));
filenames.push(path.join(__dirname, "node_modules", "protobufjs", "google", "protobuf", "type.proto"));
filenames.push(path.join(__dirname, "node_modules", "protobufjs", "google", "protobuf", "descriptor.proto"));
filenames.push(path.join(__dirname, "depscloud_api", "v1alpha", "tracker", "tracker.proto"));
filenames.push(path.join(__dirname, "depscloud_api", "v1alpha", "extractor", "extractor.proto"));
filenames.push(path.join(__dirname, "depscloud_api", "v1alpha", "schema", "schema.proto"));
filenames.push(path.join(__dirname, "depscloud_api", "v1alpha", "deps", "deps.proto"));
filenames.push(path.join(__dirname, "depscloud_api", "v1alpha", "store", "store.proto"));

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

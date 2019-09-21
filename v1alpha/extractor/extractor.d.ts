import {ChannelCredentials, Client, ServerUnaryCall, ServiceDefinition} from "grpc";
import {DependencyManagementFile} from "./deps";

export interface MatchRequest {
    separator: string;
    paths: string[];
}

export interface MatchResponse {
    matchedPaths: string[];
}

export interface ExtractRequest {
    url: string;
    separator: string;
    fileContents: { [key: string]: string };
}

export interface ExtractResponse {
    managementFiles: DependencyManagementFile[];
}

export interface IDependencyExtractor {
    match(call: ServerUnaryCall<MatchRequest>, callback: (error: Error, response: MatchResponse) => void): void;
    extract(call: ServerUnaryCall<ExtractRequest>, callback: (error: Error, response: ExtractResponse) => void): void;
}

export class DependencyExtractor extends Client {
    public static service: ServiceDefinition<IDependencyExtractor>;

    constructor(address: string, credentials: ChannelCredentials, options?: object);

    public match(request: MatchRequest, callback: (error: Error, response: MatchResponse) => void): void;
    public extract(request: ExtractRequest, callback: (error: Error, response: ExtractResponse) => void): void;
}

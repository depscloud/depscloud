import {
    ExtractRequest, ExtractResponse, MatchRequest, MatchResponse,
} from "@depscloud/api/v1alpha/extractor";
import {ServerUnaryCall} from "@grpc/grpc-js";

export default interface AsyncDependencyExtractor {
    match(request: ServerUnaryCall<MatchRequest, MatchResponse>): Promise<MatchResponse>;

    extract(request: ServerUnaryCall<ExtractRequest, ExtractResponse>): Promise<ExtractResponse>;
}

import {
    ExtractRequest, ExtractResponse, MatchRequest, MatchResponse,
} from "@depscloud/api/v1alpha/extractor";
import {ServerUnaryCall} from "@grpc/grpc-js";

export default interface AsyncDependencyExtractor {
    match(call: ServerUnaryCall<MatchRequest, MatchResponse>): Promise<MatchResponse>;

    extract(call: ServerUnaryCall<ExtractRequest, ExtractResponse>): Promise<ExtractResponse>;
}

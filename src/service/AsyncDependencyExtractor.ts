import {
    ExtractRequest, ExtractResponse, MatchRequest, MatchResponse,
} from "@deps-cloud/api/v1alpha/extractor/extractor";
import {ServerUnaryCall} from "grpc";

export default interface AsyncDependencyExtractor {
    match(request: ServerUnaryCall<MatchRequest>): Promise<MatchResponse>;

    extract(request: ServerUnaryCall<ExtractRequest>): Promise<ExtractResponse>;
}

import {ServerUnaryCall} from "grpc";
import {ExtractRequest, ExtractResponse, MatchRequest, MatchResponse} from "../../api/extractor";

export default interface AsyncDependencyExtractor {
    match(request: ServerUnaryCall<MatchRequest>): Promise<MatchResponse>;

    extract(request: ServerUnaryCall<ExtractRequest>): Promise<ExtractResponse>;
}

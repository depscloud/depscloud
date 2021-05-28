import {
    ExtractRequest, ExtractResponse, MatchRequest, MatchResponse,
} from "@depscloud/api/v1beta";
import {ServerUnaryCall} from "@grpc/grpc-js";

export default interface AsyncManifestExtractionService {
    match(call: ServerUnaryCall<MatchRequest, MatchResponse>): Promise<MatchResponse>;

    extract(call: ServerUnaryCall<ExtractRequest, ExtractResponse>): Promise<ExtractResponse>;
}

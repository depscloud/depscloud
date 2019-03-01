import {ServerUnaryCall} from "grpc";
import {ExtractRequest, ExtractResponse} from "../../api/extractor";

export default interface AsyncDependencyExtractor {
    extract(request: ServerUnaryCall<ExtractRequest>): Promise<ExtractResponse>;
}

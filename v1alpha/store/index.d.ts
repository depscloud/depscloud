import {ChannelCredentials, Client, ServerUnaryCall, ServiceDefinition} from "grpc";

type GraphItemEncoding = number;

// export const GraphItemEncoding = {
//     RAW: 0,
//     JSON: 1,
// }

export interface GraphItem {
    graphItemType: string;
    // bytes k1 = 2;
    // bytes k2 = 3;
    encoding: GraphItemEncoding;
    // bytes graphItemData = 5;
}

export interface GraphItemPair {
    edge: GraphItem;
    node: GraphItem;
}

export interface PutRequest {
    items: Array<GraphItem>;
}

export interface PutResponse {}

export interface DeleteRequest {
    items: Array<GraphItem>;
}

export interface DeleteResponse {}

export interface ListRequest {
    page: number;
    count: number;
    type: string;
}

export interface ListResponse {
    items: Array<GraphItem>;
}

export interface FindRequest {
    // bytes key = 1;
    edgeTypes: Array<string>;
}

export interface FindResponse {
    pairs: Array<GraphItemPair>;
}

export interface IGraphStore {
    put(call: ServerUnaryCall<PutRequest>, callback: (error: Error, response: PutResponse) => void): void;
    delete(call: ServerUnaryCall<DeleteRequest>, callback: (error: Error, response: DeleteResponse) => void): void;

    list(call: ServerUnaryCall<ListRequest>, callback: (error: Error, response: ListResponse) => void): void;
    findUpstream(call: ServerUnaryCall<FindRequest>, callback: (error: Error, response: FindResponse) => void): void;
    findDownstream(call: ServerUnaryCall<FindRequest>, callback: (error: Error, response: FindResponse) => void): void;
}

export class GraphStore extends Client {
    public static service: ServiceDefinition<IGraphStore>;

    constructor(address: string, credentials: ChannelCredentials, options?: object);

    public put(request: PutRequest, callback: (error: Error, response: PutResponse) => void): void;
    public delete(request: DeleteRequest, callback: (error: Error, response: DeleteResponse) => void): void;
    public list(request: ListRequest, callback: (error: Error, response: ListResponse) => void): void;
    public findUpstream(request: FindRequest, callback: (error: Error, response: FindResponse) => void): void;
    public findDownstream(request: FindRequest, callback: (error: Error, response: FindResponse) => void): void;
}

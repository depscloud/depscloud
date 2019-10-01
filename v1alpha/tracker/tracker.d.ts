import {ChannelCredentials, Client, ServerUnaryCall, Server, ServiceDefinition, ServerWritableStream, ClientReadableStream} from "grpc";
import {DependencyManagementFile} from "../deps/deps";
import {Source,Depends,Module,Manages} from "../schema/schema";

// Server side streams don't have response types so we type it

export interface TypedWriteable<T> {
    write(response: T, callback?: (error: Error) => void): void;
    end(): void;
}

export type ServerStreamCall<ReqType, RespType> = ServerWritableStream<ReqType>|TypedWriteable<RespType>;

// begin proto

export interface SourceRequest {
    source: Source;
    managementFiles: Array<DependencyManagementFile>;
}

export interface ListRequest {
    page: number;
    count: number;
}

export interface TrackResponse {
    tracking: boolean;
}

export interface ListSourceResponse {
    page: number;
    count: number;
    sources: Array<Source>;
}

export interface ISourceService {
    list(call: ServerUnaryCall<ListRequest>, callback: (error: Error, response: ListSourceResponse) => void): void;
    track(call: ServerUnaryCall<SourceRequest>, callback: (error: Error, response: TrackResponse) => void): void;
}

export class SourceService extends Client {
    public static service: ServiceDefinition<ISourceService>;

    constructor(address: string, credentials: ChannelCredentials, options?: object);

    public list(request: ListRequest, callback: (error: Error, response: ListSourceResponse) => void): void;
    public track(request: SourceRequest, callback: (error: Error, response: TrackResponse) => void): void;
}

export interface ListModuleResponse {
    page: number;
    count: number;

    modules: Array<Module>;
}

export interface ManagedSource {
    source: Source;
    manages: Manages;
}

export interface ManagedModule {
    manages: Manages;
    module: Module;
}

export interface IModuleService {
    list(call: ServerUnaryCall<ListRequest>, callback: (error: Error, response: ListModuleResponse) => void): void;
    getSource(call: ServerStreamCall<Module, ManagedSource>): void;
    getManaged(call: ServerStreamCall<Source, ManagedModule>): void;
}

export class ModuleService extends Client {
    public static service: ServiceDefinition<IModuleService>;

    constructor(address: string, credentials: ChannelCredentials, options?: object);

    public list(request: ListRequest, callback: (error: Error, response: ListModuleResponse) => void): void;
    public getSource(request: Module): ClientReadableStream<ManagedSource>;
    public getManaged(request: Source): ClientReadableStream<ManagedModule>;
}

export interface DependencyRequest {
    language: string;

    organization: string;
    module: string;
}

export interface Dependency {
    depends: Depends;
    module: Module;
}

export interface IDependencyService {
    getDependents(call: ServerStreamCall<DependencyRequest, Dependency>): void;
    getDependencies(call: ServerStreamCall<DependencyRequest, Dependency>): void;
}

export class DependencyService extends Client {
    public static service: ServiceDefinition<IDependencyService>;

    constructor(address: string, credentials: ChannelCredentials, options?: object);

    public getDependents(request: DependencyRequest): ClientReadableStream<Dependency>;
    public getDependencies(request: DependencyRequest): ClientReadableStream<Dependency>;
}

export interface TopologyTier {
    tier: Array<Dependency>;
}

export interface ITopologyService {
    getDependentsTopology(call: ServerStreamCall<DependencyRequest, Dependency>): void;
    getDependentsTopologyTiered(call: ServerStreamCall<DependencyRequest, TopologyTier>): void;
    getDependenciesTopology(call: ServerStreamCall<DependencyRequest, Dependency>): void;
    getDependenciesTopologyTiered(call: ServerStreamCall<DependencyRequest, TopologyTier>): void;
}

export class TopologyService extends Client {
    public static service: ServiceDefinition<ITopologyService>;

    constructor(address: string, credentials: ChannelCredentials, options?: object);

    public getDependentsTopology(request: DependencyRequest): ClientReadableStream<Dependency>;
    public getDependentsTopologyTiered(request: DependencyRequest): ClientReadableStream<TopologyTier>;
    public getDependenciesTopology(request: DependencyRequest): ClientReadableStream<Dependency>;
    public getDependenciesTopologyTiered(request: DependencyRequest): ClientReadableStream<TopologyTier>;
}

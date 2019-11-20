import {ChannelCredentials, Client, ServerUnaryCall, ServiceDefinition} from "grpc";
import {DependencyManagementFile} from "../deps/deps";
import {Source,Depends,Module,Manages} from "../schema/schema";

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

export interface ListSourcesResponse {
    sources: ManagedSource[];
}

export interface ListManagedResponse {
    modules: ManagedModule[];
}

export interface IModuleService {
    list(call: ServerUnaryCall<ListRequest>, callback: (error: Error, response: ListModuleResponse) => void): void;
    listSource(call: ServerUnaryCall<Module>, callback: (error: Error, response: ListSourcesResponse) => void): void;
    listManaged(call: ServerUnaryCall<Source>, callback: (error: Error, response: ListManagedResponse) => void): void;
}

export class ModuleService extends Client {
    public static service: ServiceDefinition<IModuleService>;

    constructor(address: string, credentials: ChannelCredentials, options?: object);

    public list(request: ListRequest, callback: (error: Error, response: ListModuleResponse) => void): void;
    public listSource(request: Module, callback: (error: Error, response: ListSourcesResponse) => void): void;
    public listManaged(request: Source, callback: (error: Error, response: ListManagedResponse) => void): void;
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

export interface ListDependentsResponse {
    dependents: Dependency[];
}

export interface ListDependenciesResponse {
    dependencies: Dependency[];
}

export interface IDependencyService {
    listDependents(call: ServerUnaryCall<DependencyRequest>, callback: (error: Error, response: ListDependentsResponse) => void): void;
    listDependencies(call: ServerUnaryCall<DependencyRequest>, callback: (error: Error, response: ListDependenciesResponse) => void): void;
}

export class DependencyService extends Client {
    public static service: ServiceDefinition<IDependencyService>;

    constructor(address: string, credentials: ChannelCredentials, options?: object);

    public listDependents(request: DependencyRequest, callback: (error: Error, response: ListDependentsResponse) => void): void;
    public listDependencies(request: DependencyRequest, callback: (error: Error, response: ListDependenciesResponse) => void): void;
}

export interface TopologyTier {
    tier: Array<Dependency>;
}

export interface ListDependentsTieredResponse {
    dependents: TopologyTier[];
}

export interface ListDependenciesTieredResponse {
    dependencies: TopologyTier[];
}

export interface ITopologyService {
    listDependentsTopology(call: ServerUnaryCall<DependencyRequest>, callback: (error: Error, response: ListDependentsResponse) => void): void;
    listDependentsTopologyTiered(call: ServerUnaryCall<DependencyRequest>, callback: (error: Error, response: ListDependentsTieredResponse) => void): void;
    listDependenciesTopology(call: ServerUnaryCall<DependencyRequest>, callback: (error: Error, response: ListDependenciesResponse) => void): void;
    listDependenciesTopologyTiered(call: ServerUnaryCall<DependencyRequest>, callback: (error: Error, response: ListDependenciesTieredResponse) => void): void;
}

export class TopologyService extends Client {
    public static service: ServiceDefinition<ITopologyService>;

    constructor(address: string, credentials: ChannelCredentials, options?: object);

    public listDependentsTopology(request: DependencyRequest, callback: (error: Error, response: ListDependentsResponse) => void): void;
    public listDependentsTopologyTiered(request: DependencyRequest, callback: (error: Error, response: ListDependentsTieredResponse) => void): void;
    public listDependenciesTopology(request: DependencyRequest, callback: (error: Error, response: ListDependenciesResponse) => void): void;
    public listDependenciesTopologyTiered(request: DependencyRequest, callback: (error: Error, response: ListDependenciesTieredResponse) => void): void;
}

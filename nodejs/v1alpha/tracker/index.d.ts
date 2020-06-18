import {ChannelCredentials, Client, ServerUnaryCall, ServiceDefinition} from "@grpc/grpc-js";
import {DependencyManagementFile} from "../deps";
import {Source,Depends,Module,Manages} from "../schema";

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
    list(call: ServerUnaryCall<ListRequest, ListSourceResponse>, callback: (error: Error, response: ListSourceResponse) => void): void;
    track(call: ServerUnaryCall<SourceRequest, TrackResponse>, callback: (error: Error, response: TrackResponse) => void): void;
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
    list(call: ServerUnaryCall<ListRequest, ListModuleResponse>, callback: (error: Error, response: ListModuleResponse) => void): void;
    listSource(call: ServerUnaryCall<Module, ListSourceResponse>, callback: (error: Error, response: ListSourcesResponse) => void): void;
    listManaged(call: ServerUnaryCall<Source, ListManagedResponse>, callback: (error: Error, response: ListManagedResponse) => void): void;
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
    listDependents(call: ServerUnaryCall<DependencyRequest, ListDependentsResponse>, callback: (error: Error, response: ListDependentsResponse) => void): void;
    listDependencies(call: ServerUnaryCall<DependencyRequest, ListDependenciesResponse>, callback: (error: Error, response: ListDependenciesResponse) => void): void;
}

export class DependencyService extends Client {
    public static service: ServiceDefinition<IDependencyService>;

    constructor(address: string, credentials: ChannelCredentials, options?: object);

    public listDependents(request: DependencyRequest, callback: (error: Error, response: ListDependentsResponse) => void): void;
    public listDependencies(request: DependencyRequest, callback: (error: Error, response: ListDependenciesResponse) => void): void;
}

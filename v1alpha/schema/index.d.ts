export interface Source {
    url: string;
}

export interface Manages {
    language: string;
    system: string;
    version: string;
}

export interface Module {
    language: string;
    organization: string;
    module: string;
}

export interface Depends {
    language: string;
    version_constraint: string;
    scopes: string[];
}

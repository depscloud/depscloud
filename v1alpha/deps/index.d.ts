export interface Dependency {
    organization: string;
    module: string;
    versionConstraint: string;
    scopes: string[];
}

export interface DependencyManagementFile {
    language: string;
    system: string;
    sourceUrl: string;

    organization: string;
    module: string;
    version: string;
    dependencies: Dependency[];
}

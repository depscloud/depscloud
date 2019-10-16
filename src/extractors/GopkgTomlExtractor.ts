import {Dependency, DependencyManagementFile} from "@deps-cloud/api/v1alpha/deps/deps";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import inferImportPath from "./goutils/inferImportPath";
import parseImportPath from "./goutils/parseImportPath";
import Languages from "./Languages";

interface Constraint {
    name: string;
    version: string;
    branch: string;
    revision: string;
}

function transformConstraints(data: Constraint[], scope: string): Dependency[] {
    return (data || []).map(({ name, version, branch, revision }) => {
        const { organization, module } = parseImportPath(name);

        let versionConstraint = version;
        if (branch) {
            versionConstraint = branch;
        } else if (revision) {
            versionConstraint = revision;
        }

        return {
            organization,
            module,
            versionConstraint,
            scopes: [ scope ],
        };
    });
}

function transformSimple(data: string[], versionConstraint: string, scope: string): Dependency[] {
    return (data || []).map((name) => {
        const { organization, module } = parseImportPath(name);

        return {
            organization,
            module,
            versionConstraint,
            scopes: [ scope ],
        };
    });
}

export default class GopkgTomlExtractor implements Extractor {
    public requires(): string[] {
        return [ "Gopkg.toml" ];
    }

    public async extract(url: string, files: { [p: string]: ExtractorFile }): Promise<DependencyManagementFile> {
        const { organization, module } = parseImportPath(inferImportPath(url));

        const toml = files["Gopkg.toml"].toml();

        const dependencies: Dependency[] = [];
        dependencies.push(...transformConstraints(toml.constraint, "constraint"));
        dependencies.push(...transformConstraints(toml.override, "override"));
        dependencies.push(...transformSimple(toml.ignored, "*", "ignored"));
        dependencies.push(...transformSimple(toml.required, "*", "required"));

        return {
            language: Languages.GO,
            system: "gopkg",
            organization,
            module,
            version: "latest",
            dependencies,
        };
    }
}

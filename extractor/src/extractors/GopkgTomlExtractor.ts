import {Dependency, DependencyManagementFile} from "@depscloud/api/v1alpha/deps";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import inferImportPath from "./goutils/inferImportPath";
import parseImportPath from "./goutils/parseImportPath";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";

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
            name,
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
            name,
        };
    });
}

export default class GopkgTomlExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/Gopkg.toml",
            ],
            excludes: [
                "**/vendor/**"
            ],
        };
    }

    public requires(): string[] {
        return [ "Gopkg.toml" ];
    }

    public async extract(url: string, files: { [p: string]: ExtractorFile }): Promise<DependencyManagementFile> {
        const name = inferImportPath(url);
        const { organization, module } = parseImportPath(name);

        const toml = files["Gopkg.toml"].toml();

        const dependencies: Dependency[] = [];
        dependencies.push(...transformConstraints(toml.constraint, "constraint"));
        dependencies.push(...transformConstraints(toml.override, "override"));
        dependencies.push(...transformSimple(toml.ignored, "*", "ignored"));
        dependencies.push(...transformSimple(toml.required, "*", "required"));

        return {
            language: Languages.GO,
            system: "gopkg",
            sourceUrl: "",
            organization,
            module,
            version: "latest",
            dependencies,
            name: name,
        };
    }
}

import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import inferGoImportPath from "./utils/inferGoImportPath";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";
import {ManifestDependency, ManifestFile} from "@depscloud/api/v1beta";

interface Constraint {
    name: string;
    version: string;
    branch: string;
    revision: string;
}

function transformConstraints(data: Constraint[], scope: string): ManifestDependency[] {
    return (data || []).map(({ name, version, branch, revision }) => {
        let versionConstraint = version;
        if (branch) {
            versionConstraint = branch;
        } else if (revision) {
            versionConstraint = revision;
        }

        return {
            name,
            versionConstraint,
            scopes: [ scope ],
        };
    });
}

function transformSimple(data: string[], versionConstraint: string, scope: string): ManifestDependency[] {
    return (data || []).map((name) => {
        return {
            name,
            versionConstraint,
            scopes: [ scope ],
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
                "**/vendor/**",
                "**/testdata/**",
            ],
        };
    }

    public requires(): string[] {
        return [ "Gopkg.toml" ];
    }

    public async extract(url: string, files: { [p: string]: ExtractorFile }): Promise<ManifestFile> {
        const name = inferGoImportPath(url);

        const toml = files["Gopkg.toml"].toml();

        const dependencies: ManifestDependency[] = [];
        dependencies.push(...transformConstraints(toml.constraint, "constraint"));
        dependencies.push(...transformConstraints(toml.override, "override"));
        dependencies.push(...transformSimple(toml.ignored, "*", "ignored"));
        dependencies.push(...transformSimple(toml.required, "*", "required"));

        return {
            language: Languages.GO,
            system: "gopkg",
            sourceUrl: "",
            name: name,
            version: "latest",
            dependencies,
        };
    }
}

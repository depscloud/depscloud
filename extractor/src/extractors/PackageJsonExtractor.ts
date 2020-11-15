import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";
import {ManifestDependency, ManifestFile} from "@depscloud/api/v1beta";

function extract(dependencyHash: any, scope: string): ManifestDependency[] {
    return Object.keys(dependencyHash)
        .map((dependency) => {
            return {
                name: dependency,
                versionConstraint: dependencyHash[dependency],
                scopes: scope ? [ scope ] : [],
            };
        });
}

export default class PackageJsonExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/package.json",
            ],
            excludes: [
                "**/node_modules/**",
                "**/testdata/**",
            ],
        };
    }

    public requires(): string[] {
        return [ "package.json" ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<ManifestFile> {
        const {
            name,
            version,
            repository,
            dependencies,
            devDependencies,
            peerDependencies,
            // bundledDependencies,
            optionalDependencies,
        } = files["package.json"].json();

        let allDependencies = extract((dependencies || {}), "");
        allDependencies = allDependencies.concat(extract((devDependencies || {}), "dev"));
        allDependencies = allDependencies.concat(extract((peerDependencies || {}), "peer"));
        // deps = deps.concat(extract((bundledDependencies || {}), "bundled"));
        allDependencies = allDependencies.concat(extract((optionalDependencies || {}), "optional"));

        let sourceUrl = repository;
        if (typeof repository === "object") {
            sourceUrl = repository.url;
        }

        return {
            language: Languages.NODE,
            system: "npm",
            sourceUrl,
            name,
            version,
            dependencies: allDependencies,
        };
    }
}

import {ManifestDependency, ManifestFile} from "@depscloud/api/v1beta";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";

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

export default class BowerJsonExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/bower.json",
            ],
            excludes: [
                "**/public/**",
                "**/testdata/**",
            ],
        };
    }

    public requires(): string[] {
        return [ "bower.json" ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<ManifestFile> {
        const {
            name,
            version,
            repository,
            dependencies,
            devDependencies,
            bundledDependencies,
        } = files["bower.json"].json();

        let allDependencies = extract((dependencies || {}), "");
        allDependencies = allDependencies.concat(extract((devDependencies || {}), "dev"));
        allDependencies = allDependencies.concat(extract((bundledDependencies || {}), "bundled"));

        let sourceUrl = repository;
        if (typeof repository === "object") {
            sourceUrl = repository.url;
        }

        return {
            language: Languages.JAVASCRIPT,
            system: "bower",
            sourceUrl,
            version,
            dependencies: allDependencies,
            name,
        };
    }
}

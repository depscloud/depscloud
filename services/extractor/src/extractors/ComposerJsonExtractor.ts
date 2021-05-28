import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";
import {ManifestDependency, ManifestFile} from "@depscloud/api/v1beta";

function processRequires(require: { [key: string]: string }, scope?: string): ManifestDependency[] {
    return Object.keys(require)
        .map((key) => ({ key, value: require[key] }))
        .map(({ key, value }) => {
            return {
                name: key,
                versionConstraint: value,
                scopes: scope ? [ scope ] : [],
            };
        });
}

export default class ComposerJsonExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/composer.json",
            ],
            excludes: [
                "**/vendor/**",
                "**/testdata/**",
            ],
        };
    }

    public requires(): string[] {
        return [ "composer.json" ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<ManifestFile> {
        const {
            name,
            version,
            repositories,
            require,
            "require-dev": requireDev,
        } = files["composer.json"].json();

        let dependencies = (repositories || [])
            .filter((repo) => repo.type === "package")
            .map((repo) => {
                return {
                    name: repo.package.name,
                    versionConstraint: repo.package.version,
                    scopes: ["repositories"],
                };
            });

        dependencies = dependencies.concat(processRequires(require || {}));
        dependencies = dependencies.concat(processRequires(requireDev || {}, "dev"));

        return {
            language: Languages.PHP,
            system: "composer",
            sourceUrl: "",
            name,
            version,
            dependencies,
        };
    }
}

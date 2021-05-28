import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";
import {ManifestDependency, ManifestFile} from "@depscloud/api/v1beta";

const scopes = [ "direct" ];

export default class CargoTomlExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/Cargo.toml",
            ],
            excludes: [
                "**/vendor/**",
                "**/testdata/**",
            ],
        };
    }

    public requires(): string[] {
        return [ "Cargo.toml" ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<ManifestFile> {
        const toml = files["Cargo.toml"].toml();

        const dependencies: ManifestDependency[] = Object.keys(toml.dependencies)
            .map((name) => {
                const val = toml.dependencies[name];

                if (val instanceof String) {
                    return {
                        name,
                        versionConstraint: val,
                        scopes,
                    };
                } else {
                    return {
                        name,
                        versionConstraint: val.branch,
                        scopes,
                    };
                }
            });

        return {
            language: Languages.RUST,
            system: "cargo",
            sourceUrl: "",
            name: toml.package.name,
            version: toml.package.version,
            dependencies,
        };
    }
}

import {Dependency, DependencyManagementFile} from "@depscloud/api/v1alpha/deps";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Globals from "./Globals";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";

const organization = Globals.ORGANIZATION;
const scopes = [ "direct" ];

export default class CargoTomlExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/Cargo.toml",
            ],
            excludes: [
                "**/vendor/**"
            ],
        };
    }

    public requires(): string[] {
        return [ "Cargo.toml" ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<DependencyManagementFile> {
        const toml = files["Cargo.toml"].toml();

        const dependencies: Dependency[] = Object.keys(toml.dependencies)
            .map((name) => {
                const val = toml.dependencies[name];

                if (val instanceof String) {
                    return {
                        organization,
                        module: name,
                        versionConstraint: val,
                        scopes,
                        name,
                    };
                } else {
                    return {
                        organization,
                        module: name,
                        versionConstraint: val.branch,
                        scopes,
                        name,
                    };
                }
            });

        return {
            language: Languages.RUST,
            system: "cargo",
            sourceUrl: "",
            organization,
            module: toml.package.name,
            version: toml.package.version,
            dependencies,
            name: toml.package.name,
        };
    }
}

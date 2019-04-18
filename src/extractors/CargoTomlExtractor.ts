import {Dependency, DependencyManagementFile} from "../../api/deps";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";

const organization = "__global__";
const scopes = [ "direct" ];

export default class CargoTomlExtractor implements Extractor {
    public async extract(files: { [p: string]: ExtractorFile }): Promise<DependencyManagementFile> {
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
                    };
                } else {
                    return {
                        organization,
                        module: name,
                        versionConstraint: val.branch,
                        scopes,
                    };
                }
            });

        return {
            language: "rust",
            system: "cargo",
            organization,
            module: toml.package.name,
            version: toml.package.version,
            dependencies,
        };
    }

    public requires(): string[] {
        return [ "Cargo.toml" ];
    }
}

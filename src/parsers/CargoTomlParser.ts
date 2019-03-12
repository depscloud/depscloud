import {Dependency, DependencyManagementFile} from "../../api/deps";
import {TomlParser} from "./Parser";

const organization = "__global__";
const scopes = [ "direct" ];

export default class CargoTomlParser extends TomlParser {
    public pathMatch(path: string): boolean {
        return super.pathMatch(path);
    }

    public parseToml(toml: any): DependencyManagementFile {

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
}

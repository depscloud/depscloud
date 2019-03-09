import {Dependency, DependencyManagementFile} from "../../api/deps";
import parseImportPath from "./goutils/parseImportPath";
import {TomlParser} from "./Parser";

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

export default class GopkgTomlParser extends TomlParser {
    public pathMatch(path: string): boolean {
        return path.endsWith("Gopkg.toml");
    }

    public parseToml(toml: any): DependencyManagementFile {
        const dependencies: Dependency[] = [];
        dependencies.push(...transformConstraints(toml.constraint, "constraint"));
        dependencies.push(...transformConstraints(toml.override, "override"));
        dependencies.push(...transformSimple(toml.ignored, "*", "ignored"));
        dependencies.push(...transformSimple(toml.required, "*", "required"));

        return {
            language: "go",
            system: "gopkg",
            organization: "",
            module: "",
            version: "",
            dependencies,
        };
    }
}

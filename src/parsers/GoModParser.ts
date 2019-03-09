import {Dependency, DependencyManagementFile} from "../../api/deps";
import parseImportPath from "./goutils/parseImportPath";
import {IParser} from "./Parser";

export default class GoModParser implements IParser {
    public pathMatch(path: string): boolean {
        return path.endsWith("go.mod");
    }

    public parse(path: string, content: string): DependencyManagementFile {
        const lines = content.split(/\n+/g);

        let id = null;
        const dependencies = [];

        for (let i = 0; i < lines.length; i++) {
            let line = lines[i].trim();
            if (line.length === 0) {
                continue; // empty line
            }

            const parts = line.split(/\s+/);
            const directive = parts[0];

            switch (directive) {
                case "module":
                    id = parseImportPath(parts[1]);
                    break;

                case "require":
                    i++;    // requires on subsequent lines
                    for (; i < lines.length && lines[i] !== ")"; i++) {
                        line = lines[i].trim();
                        if (line.length === 0) {
                            continue; // empty line
                        }

                        const requireParts = line.split(/\s+/);
                        const scopes = [];
                        if (requireParts.length > 2) {
                            scopes.push("indirect");
                        } else {
                            scopes.push("direct");
                        }

                        const { organization, module } = parseImportPath(requireParts[0]);
                        const dep: Dependency = {
                            organization,
                            module,
                            versionConstraint: requireParts[1],
                            scopes,
                        };

                        dependencies.push(dep);
                    }
                    break;

                case "replace":
                    i++;    // replace on subsequent lines
                    for (; i < lines.length && lines[i] !== ")"; i++) {
                        // intentionally empty
                    }
                    break;

                case "exclude":
                    i++;    // exclude on subsequent lines
                    for (; i < lines.length && lines[i] !== ")"; i++) {
                        // intentionally empty
                    }
                    break;

                default:
                    throw new Error(`parse error: unsupported directive: ${directive}`);
            }
        }

        if (id == null) {
            throw new Error(`parse error: no module present`);
        }

        return {
            language: "go",
            system: "vgo",
            organization: id.organization,
            module: id.module,
            version: "latest",
            dependencies,
        };
    }
}

import {Dependency, DependencyManagementFile} from "../../api/deps";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import parseImportPath from "./goutils/parseImportPath";

export default class GoModExtractor implements Extractor {
    public requires(): string[] {
        return [ "go.mod" ];
    }

    public async extract(files: { [p: string]: ExtractorFile }): Promise<DependencyManagementFile> {
        const content = files["go.mod"].raw();

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
                case "go":
                    break;

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
                    if (parts.length > 2) {
                        // inline replace, intentionally empty
                    } else {
                        i++;    // replace on subsequent lines
                        for (; i < lines.length && lines[i] !== ")"; i++) {
                            // intentionally empty
                        }
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
            language: "golang",
            system: "vgo",
            organization: id.organization,
            module: id.module,
            version: "latest",
            dependencies,
        };
    }
}

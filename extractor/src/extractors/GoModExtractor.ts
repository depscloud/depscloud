import {Dependency, DependencyManagementFile} from "@depscloud/api/v1alpha/deps";
import {getLogger} from "log4js";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import parseImportPath from "./goutils/parseImportPath";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";

const logger = getLogger();

export default class GoModExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/go.mod",
            ],
            excludes: [
                "**/vendor/**"
            ],
        };
    }

    public requires(): string[] {
        return [ "go.mod" ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<DependencyManagementFile> {
        const content = files["go.mod"].raw();

        const lines = content.split(/\n+/g).map(i => i.trim());

        let id = null;
        const dependencies = [];
        let name = null;

        for (let i = 0; i < lines.length; i++) {
            if (lines[i].length === 0) {
                continue; // empty line
            }

            const parts = lines[i].split(/\s+/);
            const directive = parts[0];

            switch (directive) {
                case "go":
                    break;

                case "module":
                    name = parts[1];
                    id = parseImportPath(name);
                    break;

                case "require":
                    i++;    // requires on subsequent lines
                    for (; i < lines.length && lines[i] !== ")"; i++) {
                        const line = lines[i];
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
                            name: requireParts[0],
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
                
                case "retract":
                case "exclude":
                    i++;    // exclude on subsequent lines
                    for (; i < lines.length && lines[i] !== ")"; i++) {
                        // intentionally empty
                    }
                    break;
                
                case "//":
                    break;
                
                default:
                    logger.debug(`parse error: unsupported directive: ${directive}`);
            }
        }

        if (id == null) {
            throw new Error("parse error: no module present");
        }

        return {
            language: Languages.GO,
            system: "vgo",
            sourceUrl: "",
            organization: id.organization,
            module: id.module,
            version: "latest",
            dependencies,
            name,
        };
    }
}

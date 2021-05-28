import {getLogger} from "log4js";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";
import {ManifestDependency, ManifestFile} from "@depscloud/api/v1beta";

const logger = getLogger();

export default class GoModExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/go.mod",
            ],
            excludes: [
                "**/vendor/**",
                "**/testdata/**",
            ],
        };
    }

    public requires(): string[] {
        return [ "go.mod" ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<ManifestFile> {
        const content = files["go.mod"].raw();

        const lines = content.split(/\n+/g).map(i => i.trim());

        const dependencies: Array<ManifestDependency> = [];
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

                        const dep: ManifestDependency = {
                            name: requireParts[0],
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
                    logger.debug("unsupported directive", {
                        directive,
                    });
            }
        }

        return {
            language: Languages.GO,
            system: "vgo",
            sourceUrl: "",
            name,
            version: "latest",
            dependencies,
        };
    }
}

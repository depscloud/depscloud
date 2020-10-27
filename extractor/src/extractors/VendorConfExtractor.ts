import {Dependency, DependencyManagementFile} from "@depscloud/api/v1alpha/deps";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import parseImportPath from "./goutils/parseImportPath";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";

const fileName = "vendor.conf";
const organizationString = "organization";
const moduleString = "module";

export default class VendorConfExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/vendor.conf",
            ],
            excludes: [
                "**/vendor/**",
                "**/testdata/**",
            ],
        };
    }

    public requires(): string[] {
        return [ fileName ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<DependencyManagementFile> {
        const content = files[fileName].raw();

        const lines = content.split(/\n+/g);

        let id = {};
        id[organizationString] = null;
        id[moduleString] = null;
        let idFlag = true;
        const dependencies = [];
        let name = null;

        for (const line of lines) {
            const trimmedLine = line.trim();
            if (trimmedLine.length === 0) {
                continue; // empty line
            }

            const parts = trimmedLine.split(/\s+/);
            const directive = parts[0];
            if (directive === "#") {
                continue;
            }
            if (idFlag) {
                id = parseImportPath(directive);
                name = trimmedLine;
                idFlag = false;
                continue;
            }

            const version = parts[1];
            const { organization, module } = parseImportPath(directive);

            const scopes = [];
            if (parts[2] != null) {
                scopes.push(parts[2]);
            }

            const dependencyMap: Dependency = {
                organization,
                module,
                versionConstraint: version,
                scopes,
                name: directive,
            };

            dependencies.push(dependencyMap);
        }

        if (id[organizationString] === null || id[moduleString] === null) {
            throw new Error("parse error: no module present");
        }

        return {
            language: Languages.GO,
            system: "vendor",
            sourceUrl: "",
            organization: id[organizationString],
            module: id[moduleString],
            version: "latest",
            dependencies,
            name,
        };
    }
}

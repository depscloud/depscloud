import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";
import {ManifestDependency, ManifestFile} from "@depscloud/api/v1beta";
import inferGoImportPath from "./utils/inferGoImportPath";

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
        return [ "vendor.conf" ];
    }

    public async extract(href: string, files: { [p: string]: ExtractorFile }): Promise<ManifestFile> {
        const content = files["vendor.conf"].raw();
        const name = inferGoImportPath(href);

        const lines = content.split(/\n+/g);

        const dependencies = [];

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

            const version = parts[1];

            const scopes = [];
            if (parts[2] != null) {
                scopes.push(parts[2]);
            }

            const dependencyMap: ManifestDependency = {
                name: directive,
                versionConstraint: version,
                scopes,
            };

            dependencies.push(dependencyMap);
        }

        return {
            language: Languages.GO,
            system: "vendor",
            sourceUrl: "",
            name,
            version: "latest",
            dependencies,
        };
    }
}

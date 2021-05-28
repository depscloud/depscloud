import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";
import {ManifestFile} from "@depscloud/api/v1beta";

export default class GodepsJsonExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/Godeps.json",
            ],
            excludes: [
                "**/vendor/**",
                "**/testdata/**",
            ],
        };
    }

    public requires(): string[] {
        return [ "Godeps.json" ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<ManifestFile> {
        const {
            ImportPath,
            Deps,
        } = files["Godeps.json"].json();

        const dependencies = Deps.map(({
            ImportPath: dependencyImportPath,
            Rev: version,
        }) => {

            return {
                name: dependencyImportPath,
                versionConstraint: version,
            };
        });

        return {
            language: Languages.GO,
            system: "godeps",
            sourceUrl: "",
            name: ImportPath,
            version: "latest",
            dependencies,
        };
    }
}

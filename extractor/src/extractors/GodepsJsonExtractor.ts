import {DependencyManagementFile} from "@depscloud/api/v1alpha/deps";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import parseImportPath from "./goutils/parseImportPath";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";

export default class GodepsJsonExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/Godeps.json",
            ],
            excludes: [
                "**/vendor/**"
            ],
        };
    }

    public requires(): string[] {
        return [ "Godeps.json" ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<DependencyManagementFile> {
        const {
            ImportPath,
            Deps,
        } = files["Godeps.json"].json();

        const { organization, module } = parseImportPath(ImportPath);

        const dependencies = Deps.map(({
            ImportPath: dependencyImportPath,
            Rev: version,
        }) => {
            const {
                organization: dependencyOrganization,
                module: dependencyModule,
            } = parseImportPath(dependencyImportPath);

            return {
                organization: dependencyOrganization,
                module: dependencyModule,
                versionConstraint: version,
                name: dependencyImportPath,
            };
        });

        return {
            language: Languages.GO,
            system: "godeps",
            sourceUrl: "",
            organization,
            module,
            version: "",
            dependencies,
            name: ImportPath,
        };
    }
}

import {DependencyManagementFile} from "../../api/deps";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import parseImportPath from "./goutils/parseImportPath";

export default class GodepsJsonExtractor implements Extractor {
    public requires(): string[] {
        return [ "Godeps.json" ];
    }

    public async extract(files: { [p: string]: ExtractorFile }): Promise<DependencyManagementFile> {
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
            };
        });

        return {
            language: "golang",
            system: "godeps",
            organization,
            module,
            version: "",
            dependencies,
        };
    }
}

import {DependencyManagementFile} from "../../api/deps";
import parseImportPath from "./goutils/parseImportPath";
import {JsonParser} from "./Parser";

export default class GodepsJsonParser extends JsonParser {

    public pathMatch(path: string): boolean {
        return path.endsWith("Godeps.json") && path.indexOf("vendor") === -1;
    }

    public parseJson({
        ImportPath,
        Deps,
    }: any = {}): DependencyManagementFile {
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

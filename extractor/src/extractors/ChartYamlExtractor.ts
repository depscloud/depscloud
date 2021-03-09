import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import {ManifestDependency, ManifestFile} from "@depscloud/api/v1beta";
import MatchConfig from "../matcher/MatchConfig";
import Languages from "./Languages";

export default class ChartYamlExtractor implements Extractor {
    public async extract(url: string, files: { [p: string]: ExtractorFile }): Promise<ManifestFile> {
        const {
            name,
            version,
            dependencies: deps,
        } = files["Chart.yaml"].yaml();

        const dependencies: Array<ManifestDependency> = (deps || []).map((dep) => {
            let versionConstraint = dep.repository || "";
            if (versionConstraint !== "" && dep.version) {
                versionConstraint += "#" + dep.version
            }

            return {
                name: dep.name,
                versionConstraint,
                scopes: dep.condition ? [ dep.condition ] : [],
            }
        })

        return {
            language: Languages.HELM,
            system: Languages.HELM,
            name,
            version,
            sourceUrl: "",
            dependencies,
        }
    }

    public matchConfig(): MatchConfig {
        return {
            includes: [ "**/Chart.yaml" ],
            excludes: [],
        };
    }

    public requires(): string[] {
        return [ "Chart.yaml" ];
    }
}
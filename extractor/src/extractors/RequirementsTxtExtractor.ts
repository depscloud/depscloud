import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import MatchConfig from "../matcher/MatchConfig";
import {ManifestFile} from "@depscloud/api/v1beta";
import Languages from "./Languages";
import inferPythonName from "./utils/inferPythonName";
import {ManifestDependency} from "@depscloud/api/depscloud/api/v1beta/manifest";

const operators = [
    "~", "=", "!", ">", "<"
];

export default class RequirementsTxtExtractor implements Extractor {
    public async extract(url: string, files: { [p: string]: ExtractorFile }): Promise<ManifestFile> {
        const dependencies: Array<ManifestDependency> = [];

        const contents = files["requirements.txt"].raw();
        const lines = contents.split("\n");

        lines.forEach((line) => {
            const commentIndex = line.indexOf("#");
            if (commentIndex > -1) {
                line = line.substring(0, commentIndex);
            }

            line = line.trim()

            if (!line || line[0] === "-") {
                return
            }

            let l = line.length;
            for (const operator of operators) {
                const i = line.indexOf(operator);
                if (i > -1) {
                    l = Math.min(l, i);
                }
            }

            const name = line.substring(0, l).trim();
            const versionConstraint = line.substring(l).trim();

            dependencies.push({
                name,
                versionConstraint: (versionConstraint || "*"),
                scopes: [""],
            });
        });

        return {
            language: Languages.PYTHON,
            system: "pip",
            sourceUrl: "",
            name: inferPythonName(url),
            version: "*",
            dependencies,
        };
    }

    public matchConfig(): MatchConfig {
        return {
            includes: [ "**/requirements.txt" ],
            excludes: [],
        };
    }

    public requires(): string[] {
        return [ "requirements.txt" ];
    }
}

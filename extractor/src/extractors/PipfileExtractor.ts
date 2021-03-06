import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import MatchConfig from "../matcher/MatchConfig";
import {ManifestDependency, ManifestFile} from "@depscloud/api/v1beta";
import Languages from "./Languages";
import inferPythonName from "./utils/inferPythonName";

const transform = (scope: string, name: string, val: any): ManifestDependency => {
    let versionConstraint = "*";
    if (typeof val == "string") {
        versionConstraint = val;
    } else if (val.version) {
        versionConstraint = val.version;
    } else if (val.ref) { // git
        versionConstraint = val.ref;
    }

    return {
        name,
        versionConstraint,
        scopes: scope ? [ scope ] : [],
    }
}

export default class PipfileExtractor implements Extractor {
    public async extract(url: string, files: { [p: string]: ExtractorFile }): Promise<ManifestFile> {
        const {
            source,
            packages,
            "dev-packages": devPackages,
        } = files["Pipfile"].toml();

        let [ name, sourceUrl ] = [ "", "" ];
        if (source) {
            name = source[0].name;
            sourceUrl = source[0].url;
        } else {
            name = inferPythonName(url);
        }

        const dependencies = [];

        Object.keys(packages)
            .map((key) => transform("", key, packages[key]))
            .forEach((dep) => dependencies.push(dep));

        Object.keys(devPackages)
            .map((key) => transform("dev", key, devPackages[key]))
            .forEach((dep) => dependencies.push(dep));

        return {
            language: Languages.PYTHON,
            system: "pipfile",
            name,
            version: "*",
            sourceUrl,
            dependencies,
        };
    }

    public matchConfig(): MatchConfig {
        return {
            includes: [ "**/Pipfile" ],
            excludes: [],
        };
    }

    public requires(): string[] {
        return [ "Pipfile" ];
    }
}
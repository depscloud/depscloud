import MatchConfig from "../matcher/MatchConfig";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Languages from "./Languages";
import {ManifestDependency, ManifestFile} from "@depscloud/api/v1beta";

const extract = (dependency: any):ManifestDependency => {
    const { source, version } = dependency;
    const { git } = (source || {});
    const { subdir, remote } = (git || {});

    return {
        name: remote,
        versionConstraint: version,
        scopes: subdir ? [ subdir ] : [],
    };
}

export default class JsonnetfileJsonExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/jsonnetfile.json",
            ],
            excludes: [],
        };
    }

    public requires(): string[] {
        return [ "jsonnetfile.json" ];
    }

    public async extract(url: string, files: { [key: string]: ExtractorFile }): Promise<ManifestFile> {
        const {
            version,
            dependencies,
        } = files["jsonnetfile.json"].json();

        const deps = dependencies.map(dependency => extract(dependency));


        return {
            language: Languages.JSONNET,
            system:"jsonnet-bundler",
            dependencies: deps,
            name: url,
            sourceUrl:"",
            version
        }
    }
}

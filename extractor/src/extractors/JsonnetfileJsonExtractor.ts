import {Dependency} from "@depscloud/api/v1alpha/deps";
import MatchConfig from "../matcher/MatchConfig";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Languages from "./Languages";

const extract = (dependency: any,scopes:any) => {
    return {
        scopes:[scopes, dependency.source.git.subdir ],
        name: dependency.source.git.remote,
    };
}

export default class JsonnetfileJsonExtractor implements Extractor {
    matchConfig(): MatchConfig {
        return {
            includes: [
                "**/jsonnetfile.json",
            ],
            excludes: [],
        };
    }
    requires(): string[] {
        return [ "jsonnetfile.json" ];
    }
    public async extract(url: string, files: { [key: string]: ExtractorFile }): Promise<any> {
        const {
           dependencies,
        } = files["jsonnetfile.json"].json();

        const deps = dependencies.map(dependency => extract(dependency , "direct"));


        return {
            language: Languages.JSONNET,
            system:"jsonnet-bundler",
            dependencies: deps,
            name,
        }
    }
}
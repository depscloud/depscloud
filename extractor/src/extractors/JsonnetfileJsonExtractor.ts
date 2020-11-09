import {Dependency, DependencyManagementFile} from "@depscloud/api/v1alpha/deps";
import MatchConfig from "../matcher/MatchConfig";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Languages from "./Languages";

const extract = (dependency: any):Dependency => {
    const { source, version } = dependency;
  const { git } = (source || {});
  const { subdir, name } = (git || {});

  return {
    name,
    versionConstraint: version,
    scopes: subdir ? [ subdir ] : [],
    organization: "", // deprecated
    module: "", // deprecated
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
    public async extract(url: string, files: { [key: string]: ExtractorFile }): Promise<DependencyManagementFile> {
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
            organization:"",
            module:"", //deprecated
            sourceUrl:"",
            version
        }
    }
}
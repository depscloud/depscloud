import {Dependency, DependencyManagementFile} from "@depscloud/api/v1alpha/deps";
import MatchConfig from "../matcher/MatchConfig";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Languages from "./Languages";

const extract = (dependency: any,scopes:any):Dependency => {
    return {
        organization:"",
        module:"",
        versionConstraint: dependency.version,
        scopes:[scopes],
        name: parseName(dependency),
    };
}

//The name should be formatted as {{ remote }}[#{{ subdir }}] where subdir
const parseName = (dependency) => {
    return `${dependency.source.git.remote}[#${dependency.source.git.subdir}]`;
}

    

class JsonnetJsonExtractor implements Extractor {
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
            name,
            version,
            organization,
            dependencies,

        } = files["jsonnetfile.json"].json();

        const deps = dependencies.map(dependence => extract(dependence , "dependencies"));


        return {
            language: Languages.JAVASCRIPT,
            system:"",
            sourceUrl: "",
            organization,
            module: "",
            version,
            dependencies: deps,
            name,
        }
    }
}

export default  JsonnetJsonExtractor;
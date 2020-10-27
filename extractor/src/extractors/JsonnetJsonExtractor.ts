import {Dependency, DependencyManagementFile} from "@depscloud/api/v1alpha/deps";
import MatchConfig from "../matcher/MatchConfig";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Languages from "./Languages";

const extract = (dependency: any, scope: string) => {
    return {
        name: parseName(dependency),
        organization:"",
        module,
        versionConstraint: dependency.version,
        scopes: [ scope ],
    };
}

const parseName = (dependency) => {
    return dependency.source.git.remote +"["+"#"+ dependency.source.git.subdir+"]";
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
            version,
            dependencies,
            legacyImports,
        } = files["jsonnetfile.json"].json();

        const deps = dependencies.map(dependence => extract(dependence,""));


        return {
            language: Languages.JAVASCRIPT,
            system: "npm",
            sourceUrl: "",
            organization: "",
            module: "string",
            version,
            dependencies: deps,
            name: "",
        }
    }
}

export default  JsonnetJsonExtractor;
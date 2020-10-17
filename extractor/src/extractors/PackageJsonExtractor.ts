import {Dependency, DependencyManagementFile} from "@depscloud/api/v1alpha/deps";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Globals from "./Globals";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";

interface ID {
    organization: string;
    module: string;
}

function parseName(module: string): ID {
    let organization = Globals.ORGANIZATION;
    if (module.charAt(0) === "@") {
        const index = module.indexOf("/");
        organization = module.substr(1,  index - 1);
        module = module.substr(index + 1);
    }
    return { organization, module };
}

function extract(dependencyHash: any, scope: string): Dependency[] {
    return Object.keys(dependencyHash)
        .map((dependency) => {
            const { organization, module } = parseName(dependency);
            const versionConstraint = dependencyHash[dependency];

            return {
                organization,
                module,
                versionConstraint,
                scopes: [ scope ],
                name: dependency,
            };
        });
}

export default class PackageJsonExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/package.json",
            ],
            excludes: [
                "**/node_modules/**"
            ],
        };
    }

    public requires(): string[] {
        return [ "package.json" ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<DependencyManagementFile> {
        const {
            name,
            version,
            repository,
            dependencies,
            devDependencies,
            peerDependencies,
            // bundledDependencies,
            optionalDependencies,
        } = files["package.json"].json();

        const { organization, module } = parseName((name || ""));

        let allDependencies = extract((dependencies || {}), "");
        allDependencies = allDependencies.concat(extract((devDependencies || {}), "dev"));
        allDependencies = allDependencies.concat(extract((peerDependencies || {}), "peer"));
        // deps = deps.concat(extract((bundledDependencies || {}), "bundled"));
        allDependencies = allDependencies.concat(extract((optionalDependencies || {}), "optional"));

        let sourceUrl = repository;
        if (typeof repository === "object") {
            sourceUrl = repository.url;
        }

        return {
            language: Languages.NODE,
            system: "npm",
            sourceUrl,
            organization, module, version,
            dependencies: allDependencies,
            name,
        };
    }
}

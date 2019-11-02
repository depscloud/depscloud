import {Dependency, DependencyManagementFile} from "@deps-cloud/api/v1alpha/deps/deps";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Globals from "./Globals";
import Languages from "./Languages";

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
            };
        });
}

export default class PackageJsonExtractor implements Extractor {
    public requires(): string[] {
        return [ "package.json" ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<DependencyManagementFile> {
        const {
            name,
            version,
            dependencies,
            devDependencies,
            peerDependencies,
            bundledDependencies,
            optionalDependencies,
        } = files["package.json"].json();

        const { organization, module } = parseName((name || ""));

        let allDependencies = extract((dependencies || {}), "");
        allDependencies = allDependencies.concat(extract((devDependencies || {}), "dev"));
        allDependencies = allDependencies.concat(extract((peerDependencies || {}), "peer"));
        // deps = deps.concat(extract((bundledDependencies || {}), "bundled"));
        allDependencies = allDependencies.concat(extract((optionalDependencies || {}), "optional"));

        return {
            language: Languages.NODE,
            system: "npm",
            organization, module, version,
            dependencies: allDependencies,
        };
    }
}

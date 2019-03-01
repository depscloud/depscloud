import { Dependency, DependencyManagementFile } from "../../api/deps";
import {JsonParser} from "./Parser";

interface ID {
    organization: string;
    module: string;
}

function parseName(module: string): ID {
    let organization = "_";
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

export default class PackageJsonParser extends JsonParser {
    public pathMatch(path: string): boolean {
        return path.endsWith("package.json") && path.indexOf("node_modules") === -1;
    }

    public parseJson({
        name,
        version,
        dependencies,
        devDependencies,
        peerDependencies,
        bundledDependencies,
        optionalDependencies,
    }): DependencyManagementFile {
        const { organization, module } = parseName(name);

        let allDependencies = extract((dependencies || {}), "");
        allDependencies = allDependencies.concat(extract((devDependencies || {}), "dev"));
        allDependencies = allDependencies.concat(extract((peerDependencies || {}), "peer"));
        // deps = deps.concat(extract((bundledDependencies || {}), "bundled"));
        allDependencies = allDependencies.concat(extract((optionalDependencies || {}), "optional"));

        return {
            language: "node",
            system: "npm",
            organization, module, version,
            dependencies: allDependencies,
        };
    }

}

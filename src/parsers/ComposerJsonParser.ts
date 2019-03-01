import { Dependency, DependencyManagementFile } from "../../api/deps";
import {JsonParser} from "./Parser";

interface ID {
    organization: string;
    module: string;
}

function parseName(name: string): ID {
    const pos = name.indexOf("/");
    const organization = name.substr(0, pos);
    const module = name.substr(pos + 1);
    return { organization, module };
}

function processRequires(require: { [key: string]: string }): Dependency[] {
    return Object.keys(require)
        .map((key) => ({ key, value: require[key] }))
        .map(({ key, value }) => {
            const { organization, module } = parseName(key);

            return {
                organization,
                module,
                versionConstraint: value,
            } as Dependency;
        });
}

export default class ComposerJsonParser extends JsonParser {
    public pathMatch(path: string): boolean {
        return path.endsWith("composer.json") && path.indexOf("vendor") === -1;
    }

    public parseJson({
        name,
        version,
        repositories,
        require,
        "require-dev": requireDev,
    }): DependencyManagementFile {
        const { organization, module } = parseName(name);

        let dependencies = (repositories || [])
            .filter((repo) => repo.type === "package")
            .map((repo) => {
                const {
                    organization: dependencyOrganization,
                    module: dependencyModule,
                } = parseName(repo.package.name);

                return {
                    organization: dependencyOrganization,
                    module: dependencyModule,
                    versionConstraint: repo.package.version,
                };
            });

        dependencies = dependencies.concat(processRequires(require || {}));
        dependencies = dependencies.concat(processRequires(requireDev || {}));

        return {
            language: "php",
            system: "composer",
            organization,
            module,
            version,
            dependencies,
        };
    }
}

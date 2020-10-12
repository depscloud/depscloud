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

function parseName(name: string): ID {
    const pos = name.indexOf("/");
    if (pos === -1) {
        return { organization: Globals.ORGANIZATION, module: name };
    }

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

export default class ComposerJsonExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/composer.json",
            ],
            excludes: [
                "**/vendor/**"
            ],
        };
    }

    public requires(): string[] {
        return [ "composer.json" ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<DependencyManagementFile> {
        const {
            name,
            version,
            repositories,
            require,
            "require-dev": requireDev,
        } = files["composer.json"].json();

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
                    name: repo.package.name,
                };
            });

        dependencies = dependencies.concat(processRequires(require || {}));
        dependencies = dependencies.concat(processRequires(requireDev || {}));

        return {
            language: Languages.PHP,
            system: "composer",
            sourceUrl: "",
            organization,
            module,
            version,
            dependencies,
            name,
        };
    }
}

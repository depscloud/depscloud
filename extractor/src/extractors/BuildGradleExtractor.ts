import {parseText} from "gradle-to-js/lib/parser";

import {Dependency, DependencyManagementFile} from "@depscloud/api/v1alpha/deps";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Globals from "./Globals";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";

// infer the module name of the last segment of the git url
function inferModuleName(url: string): string {
   return url.substring(url.lastIndexOf("/"), url.length - 4);
}

export default class BuildGradleExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/build.gradle",
                "**/settings.gradle",
            ],
            excludes: [
                "**/testdata/**",
            ],
        };
    }

    public requires(): string[] {
        return [ "build.gradle", "settings.gradle" ];
    }

    public async extract(url: string, files: { [p: string]: ExtractorFile }): Promise<DependencyManagementFile> {
        const promises = this.requires()
            .map((req) => files[req].raw())
            .map((raw) => parseText(raw));

        const [
            buildGradle,
            settingsGradle,
        ] = await Promise.all(promises);

        const dependencies: { [key: string]: Dependency } = {};

        let [ organization, module ] = [ Globals.ORGANIZATION, inferModuleName(url) ];

        if (settingsGradle.rootProject) {
            if (settingsGradle.rootProject.name) {
                module = settingsGradle.rootProject.name;
            }

            if (settingsGradle.rootProject.parent) {
                const key = settingsGradle.rootProject.parent.name;
                const [ parentOrganization, parentModule, parentVersion ] = key.split(":");

                dependencies[key] = {
                    organization: parentOrganization,
                    module: parentModule,
                    versionConstraint: parentVersion,
                    scopes: [ "parent" ],
                    name: [ parentOrganization, parentModule ].join(":"),
                };
            }
        }

        if (buildGradle.group) {
            organization = buildGradle.group;
        }

        (buildGradle.dependencies || []).forEach((dep) => {
            const key = [ dep.group, dep.name, dep.version ].join(":");

            if (dependencies[key]) {
                dependencies[key].scopes.push(dep.type);
            } else {
                dependencies[key] = {
                    organization: dep.group,
                    module: dep.name,
                    versionConstraint: dep.version,
                    scopes: [ dep.type ],
                    name: [ dep.group, dep.name ].join(":"),
                };
            }
        });

        return {
            language: Languages.JAVA,
            system: "gradle",
            sourceUrl: "",
            organization,
            module,
            version: buildGradle.version,
            dependencies: Object.keys(dependencies)
                .map((k) => dependencies[k]),
            name: [ organization, module ].join(":"),
        };
    }
}

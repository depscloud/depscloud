import {parseText} from "gradle-to-js/lib/parser";

import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";
import {ManifestDependency, ManifestFile} from "@depscloud/api/v1beta";

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

    public async extract(url: string, files: { [p: string]: ExtractorFile }): Promise<ManifestFile> {
        const promises = this.requires()
            .map((req) => files[req].raw())
            .map((raw) => parseText(raw));

        const [
            buildGradle,
            settingsGradle,
        ] = await Promise.all(promises);

        const dependencies: { [key: string]: ManifestDependency } = {};
        let [ organization, module ] = [ "", "" ]

        if (settingsGradle.rootProject) {
            if (settingsGradle.rootProject.name) {
                module = settingsGradle.rootProject.name;
            }

            if (settingsGradle.rootProject.parent) {
                const key = settingsGradle.rootProject.parent.name;
                const [ parentOrganization, parentModule, parentVersion ] = key.split(":");

                dependencies[key] = {
                    name: [ parentOrganization, parentModule ].join(":"),
                    versionConstraint: parentVersion,
                    scopes: [ "parent" ],
                };
            }
        }

        if (buildGradle.group) {
            organization = buildGradle.group;
        }

        (buildGradle.dependencies || []).forEach((dep) => {
            const key = [ dep.group, dep.name ].join(":");

            if (dependencies[key]) {
                dependencies[key].scopes.push(dep.type);
            } else {
                dependencies[key] = {
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
            name: organization ? [ organization, module ].join(":") : module,
            version: buildGradle.version,
            dependencies: Object.keys(dependencies)
                .map((k) => dependencies[k]),
        };
    }
}

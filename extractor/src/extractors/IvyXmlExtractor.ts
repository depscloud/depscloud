import {Dependency, DependencyManagementFile} from "@depscloud/api/v1alpha/deps";
import cheerio = require("cheerio");
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Globals from "./Globals";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";

export default class IvyXmlExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/ivy.xml",
            ],
            excludes: [],
        };
    }

    public requires(): string[] {
        return [ "ivy.xml" ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<DependencyManagementFile> {
        const xml = files["ivy.xml"].xml();

        const infoNode: Cheerio = xml.find("ivy-module info");
        const dependencyNodes: Cheerio = xml.find("ivy-module dependencies dependency");

        const dependencies: Dependency[] = [];
        dependencyNodes.map((i, dependencyNode: CheerioElement) => {
            const dep = cheerio(dependencyNode);
            const confs = dep.attr("conf");
            const organization = dep.attr("org");
            const module = dep.attr("name");
            const versionConstraint = dep.attr("rev");

            let scopes = [];
            if (confs) {
                scopes = confs.split(";");
            }

            dependencies.push({ organization, module, versionConstraint, scopes, name: [ organization , module ].join(":") });
        });

        return {
            language: Languages.JAVA,
            system: "ivy",
            sourceUrl: "",
            organization: infoNode.attr("organisation") || Globals.ORGANIZATION,
            module: infoNode.attr("module"),
            version: infoNode.attr("revision") || null,
            dependencies,
            name: [ infoNode.attr("organisation"), infoNode.attr("module") ].filter(Boolean).join(":"),
        };
    }
}

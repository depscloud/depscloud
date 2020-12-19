import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";
import {ManifestDependency, ManifestFile} from "@depscloud/api/v1beta";
import cheerio = require("cheerio");

export default class IvyXmlExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/ivy.xml",
            ],
            excludes: [
                "**/testdata/**",
            ],
        };
    }

    public requires(): string[] {
        return [ "ivy.xml" ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<ManifestFile> {
        const xml = files["ivy.xml"].xml();

        const infoNode: cheerio.Cheerio = xml.find("ivy-module info");
        const dependencyNodes: cheerio.Cheerio = xml.find("ivy-module dependencies dependency");

        const dependencies: ManifestDependency[] = [];
        dependencyNodes.map((i, dependencyNode: cheerio.Element) => {
            const dep = cheerio(dependencyNode);
            const confs = dep.attr("conf");
            const organization = dep.attr("org");
            const module = dep.attr("name");
            const versionConstraint = dep.attr("rev");

            let scopes = [];
            if (confs) {
                scopes = confs.split(";");
            }

            dependencies.push({
                name: [ organization , module ].join(":"),
                versionConstraint,
                scopes,
            });
        });

        return {
            language: Languages.JAVA,
            system: "ivy",
            sourceUrl: "",
            name: [ infoNode.attr("organisation"), infoNode.attr("module") ].filter(Boolean).join(":"),
            version: infoNode.attr("revision") || null,
            dependencies,
        };
    }
}

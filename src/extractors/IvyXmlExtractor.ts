import cheerio = require("cheerio");
import {Dependency, DependencyManagementFile} from "../../api/deps";
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Globals from "./Globals";
import Languages from "./Languages";

export default class IvyXmlExtractor implements Extractor {
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

            dependencies.push({ organization, module, versionConstraint, scopes });
        });

        return {
            language: Languages.JAVA,
            system: "ivy",
            organization: infoNode.attr("organisation") || Globals.ORGANIZATION,
            module: infoNode.attr("module"),
            version: infoNode.attr("revision") || null,
            dependencies,
        };
    }
}

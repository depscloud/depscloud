import cheerio = require("cheerio");
import {Dependency, DependencyManagementFile} from "../../api/deps";
import {XmlParser} from "./Parser";

export default class IvyXmlParser extends XmlParser {
    public pathMatch(path: string): boolean {
        return path.endsWith("ivy.xml");
    }

    public parseXml(xml: Cheerio): DependencyManagementFile {
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
            language: "java",
            system: "ivy",
            organization: infoNode.attr("organisation"),
            module: infoNode.attr("module"),
            version: infoNode.attr("revision") || null,
            dependencies,
        };
    }
}

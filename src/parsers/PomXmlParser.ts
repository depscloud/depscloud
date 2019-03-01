import cheerio = require("cheerio");
import { Dependency, DependencyManagementFile } from "../../api/deps";
import {XmlParser} from "./Parser";

export default class PomXmlParser extends XmlParser {
    public pathMatch(path: string): boolean {
        return path.endsWith("pom.xml");
    }

    public parseXml(xml: Cheerio): DependencyManagementFile {
        const parentGroupId = xml.find("project > parent > groupId").text();
        const parentArtifactId = xml.find("project > parent > artifactId").text();
        const parentVersion = xml.find("project > parent > version").text();

        const groupId = xml.find("project > groupId").text() || parentGroupId;
        const artifactId = xml.find("project > artifactId").text();
        const version = xml.find("project > version").text();

        const dependencies: Dependency[] = [];
        if (parentGroupId && parentArtifactId && parentVersion) {
            dependencies.push({
                organization: parentGroupId,
                module: parentArtifactId,
                versionConstraint: parentVersion,
                scopes: [ "parent" ],
            });
        }

        const matched: Cheerio = xml.find("project > dependencies > dependency");
        matched.map((i, match: CheerioElement) => {
            const organization = cheerio(match).find("groupId").text();
            const module = cheerio(match).find("artifactId").text();
            const versionConstraint = cheerio(match).find("version").text();
            const scope = cheerio(match).find("scope").text();

            const scopes = [scope || "compile"];

            dependencies.push({organization, module, versionConstraint, scopes});
        });

        return {
            language: "java",
            system: "maven",
            organization: groupId,
            module: artifactId,
            version,
            dependencies,
        };
    }
}

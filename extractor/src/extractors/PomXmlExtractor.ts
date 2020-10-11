import {Dependency, DependencyManagementFile} from "@depscloud/api/v1alpha/deps";
import cheerio = require("cheerio");
import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";

export default class PomXmlExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/pom.xml",
                "**/*.pom",
            ],
            excludes: [],
        };
    }

    public requires(): string[] {
        return [ "pom.xml" ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<DependencyManagementFile> {
        const xml = files["pom.xml"].xml();

        const parentGroupId = xml.find("project > parent > groupId").text();
        const parentArtifactId = xml.find("project > parent > artifactId").text();
        const parentVersion = xml.find("project > parent > version").text();

        const groupId = xml.find("project > groupId").text() || parentGroupId;
        const artifactId = xml.find("project > artifactId").text();
        const version = xml.find("project > version").text();

        const sourceUrl = xml.find("project > scm > url").text();

        const dependencies: Dependency[] = [];
        if (parentGroupId && parentArtifactId && parentVersion) {
            dependencies.push({
                organization: parentGroupId,
                module: parentArtifactId,
                versionConstraint: parentVersion,
                scopes: [ "parent" ],
                name: [parentGroupId, parentArtifactId].join(";"),
            });
        }

        const matched: Cheerio = xml.find("project > dependencies > dependency");
        matched.map((i, match: CheerioElement) => {
            const organization = cheerio(match).find("> groupId").text();
            const module = cheerio(match).find("> artifactId").text();
            const versionConstraint = cheerio(match).find("> version").text();
            const scope = cheerio(match).find("> scope").text();

            const scopes = [scope || "compile"];

            dependencies.push({organization, module, versionConstraint, scopes, name: [groupId, artifactId].join(";")});
        });

        return {
            language: Languages.JAVA,
            system: "maven",
            sourceUrl,
            organization: groupId,
            module: artifactId,
            version,
            dependencies,
            name: [groupId, artifactId].join(";"),
        };
    }
}

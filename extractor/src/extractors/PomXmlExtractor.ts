import Extractor from "./Extractor";
import ExtractorFile from "./ExtractorFile";
import Languages from "./Languages";
import MatchConfig from "../matcher/MatchConfig";
import {ManifestDependency, ManifestFile} from "@depscloud/api/v1beta";
import cheerio = require("cheerio");

export default class PomXmlExtractor implements Extractor {
    public matchConfig(): MatchConfig {
        return {
            includes: [
                "**/pom.xml",
                "**/*.pom",
            ],
            excludes: [
                "**/testdata/**",
            ],
        };
    }

    public requires(): string[] {
        return [ "pom.xml" ];
    }

    public async extract(_: string, files: { [p: string]: ExtractorFile }): Promise<ManifestFile> {
        const xml = files["pom.xml"].xml();

        const parentGroupId = xml.find("project > parent > groupId").text();
        const parentArtifactId = xml.find("project > parent > artifactId").text();
        const parentVersion = xml.find("project > parent > version").text();

        const groupId = xml.find("project > groupId").text() || parentGroupId;
        const artifactId = xml.find("project > artifactId").text();
        const version = xml.find("project > version").text();

        const sourceUrl = xml.find("project > scm > url").text();

        const dependencies: ManifestDependency[] = [];
        if (parentGroupId && parentArtifactId && parentVersion) {
            dependencies.push({
                name: [parentGroupId, parentArtifactId].join(":"),
                versionConstraint: parentVersion,
                scopes: [ "parent" ],
            });
        }

        const matched: Cheerio = xml.find("project > dependencies > dependency");
        matched.map((i, match: CheerioElement) => {
            const depGroupId = cheerio(match).find("> groupId").text();
            const depArtifactId = cheerio(match).find("> artifactId").text();
            const versionConstraint = cheerio(match).find("> version").text();
            const scope = cheerio(match).find("> scope").text();

            const scopes = [scope || "compile"];

            dependencies.push({
                name: [depGroupId, depArtifactId].join(":"),
                versionConstraint,
                scopes,
            });
        });

        return {
            language: Languages.JAVA,
            system: "maven",
            sourceUrl,
            name: [groupId, artifactId].join(":"),
            version,
            dependencies,
        };
    }
}

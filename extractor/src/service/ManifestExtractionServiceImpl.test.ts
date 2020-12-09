import fs = require("fs");
import path = require("path");
import ExtractorRegistry from "../extractors/ExtractorRegistry";
import ManifestExtractionServiceImpl from "./ManifestExtractionServiceImpl";
import Matcher from "../matcher/Matcher";

const fsp = fs.promises;

const files = [
    "BuildGradleExtractor.test.ts",
    "BuildGradleExtractor.ts",
    "CargoTomlExtractor.test.ts",
    "CargoTomlExtractor.ts",
    "ComposerJsonExtractor.test.ts",
    "ComposerJsonExtractor.ts",
    "Extractor.ts",
    "ExtractorFile.ts",
    "ExtractorRegistry.ts",
    "GoModExtractor.test.ts",
    "GoModExtractor.ts",
    "GodepsJsonExtractor.test.ts",
    "GodepsJsonExtractor.ts",
    "GopkgTomlExtractor.test.ts",
    "GopkgTomlExtractor.ts",
    "IvyXmlExtractor.test.ts",
    "IvyXmlExtractor.ts",
    "PackageJsonExtractor.test.ts",
    "PackageJsonExtractor.ts",
    "PomXmlExtractor.test.ts",
    "PomXmlExtractor.ts",
    "Cargo.toml",
    "Godeps.json",
    "Gopkg.toml",
    "build.gradle",
    "composer.json",
    "go.mod",
    "go.sum",
    "ivy.xml",
    "package.json",
    "pom.xml",
    "settings.gradle",
];

describe("ManifestExtractionServiceImpl", () => {
    test("fullParse", async () => {
        const extractorPromises = ExtractorRegistry.known()
            .map((registry) => ExtractorRegistry.resolve(registry, null));

        const extractors = await Promise.all(extractorPromises);

        const matchersAndExtractors = extractors.map((extractor) => {
            return {
                matcher: new Matcher(extractor.matchConfig()),
                extractor,
            }
        });

        const extractorImpl = new ManifestExtractionServiceImpl(matchersAndExtractors);

        const matched = extractorImpl.matchInternal(path.sep, files);

        expect(matched.sort()).toMatchSnapshot();

        const extractorsDir = path.resolve(__dirname, "../extractors/testdata");

        const fileContents = {};
        const promises = matched.map((match) =>
            fsp.readFile(path.join(extractorsDir, match))
                .then((buf) => buf.toString())
                .then((content) => fileContents[match] = content));

        await Promise.all(promises);

        const dependencyManagementFiles = await extractorImpl.extractInternal(
            "git@github.com:depscloud/extractor.git", path.sep, fileContents);

        expect(dependencyManagementFiles).toMatchSnapshot();
    });
});

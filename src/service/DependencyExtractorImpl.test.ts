import fs = require("fs");
import path = require("path");
import ExtractorRegistry from "../extractors/ExtractorRegistry";
import DependencyExtractorImpl from "./DependencyExtractorImpl";

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
    "testdata/Cargo.toml",
    "testdata/Godeps.json",
    "testdata/Gopkg.toml",
    "testdata/build.gradle",
    "testdata/composer.json",
    "testdata/go.mod",
    "testdata/go.sum",
    "testdata/ivy.xml",
    "testdata/package.json",
    "testdata/pom.xml",
    "testdata/settings.gradle",
];

describe("DependencyExtractorImpl", () => {
    test("fullParse", async () => {
        const extractorPromises = ExtractorRegistry.known()
            .map((registry) => ExtractorRegistry.resolve(registry, null));

        const extractors = await Promise.all(extractorPromises);

        const extractorImpl = new DependencyExtractorImpl(extractors);

        const matched = extractorImpl.matchInternal(path.sep, files);

        expect(matched).toMatchSnapshot();

        const extractorsDir = path.resolve(__dirname, "../extractors");

        const fileContents = {};
        const promises = matched.map((match) =>
            fsp.readFile(path.join(extractorsDir, match))
                .then((buf) => buf.toString())
                .then((content) => fileContents[match] = content));

        await Promise.all(promises);

        const dependencyManagementFiles = await extractorImpl.extractInternal(path.sep, fileContents);

        expect(dependencyManagementFiles).toMatchSnapshot();
    });
});

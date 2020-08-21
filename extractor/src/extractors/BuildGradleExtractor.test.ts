import {readFile} from "fs";
import {promisify} from "util";
import BuildGradleExtractor from "./BuildGradleExtractor";
import ExtractorFile from "./ExtractorFile";

const readFileAsync = promisify(readFile);

describe("CargoTomlExtractor", () => {
    test("fullParse", async () => {
        const buildGradle = require.resolve("./testdata/build.gradle");
        const settingsGradle = require.resolve("./testdata/settings.gradle");

        const [
            buildGradleContent,
            settingsGradleContent,
        ] = await Promise.all([
            readFileAsync(buildGradle).then((buf) => buf.toString()),
            readFileAsync(settingsGradle).then((buf) => buf.toString()),
        ]);

        const parser = new BuildGradleExtractor();

        const actual = await parser.extract("", {
            "build.gradle": new ExtractorFile(buildGradleContent),
            "settings.gradle": new ExtractorFile(settingsGradleContent),
        });

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

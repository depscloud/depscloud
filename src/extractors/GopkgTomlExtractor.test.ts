import {readFile} from "fs";
import {promisify} from "util";
import ExtractorFile from "./ExtractorFile";
import GopkgTomlExtractor from "./GopkgTomlExtractor";

const readFileAsync = promisify(readFile);

describe("GopkgTomlExtractor", () => {
    test("fullParse", async () => {
        const gopkgToml = require.resolve("./testdata/Gopkg.toml");
        const buffer = await readFileAsync(gopkgToml);
        const content = buffer.toString();

        const parser = new GopkgTomlExtractor();

        const actual = await parser.extract(
            "git@github.com:deps-cloud/extractor.git", {
                "Gopkg.toml": new ExtractorFile(content),
            },
        );

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

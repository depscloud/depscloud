import {readFile} from "fs";
import {promisify} from "util";
import ComposerJsonExtractor from "./ComposerJsonExtractor";
import ExtractorFile from "./ExtractorFile";

const readFileAsync = promisify(readFile);

describe("ComposerJsonExtractor", () => {
    test("fullParse", async () => {
        const composerPath = require.resolve("./testdata/composer.json");
        const buffer = await readFileAsync(composerPath);
        const content = buffer.toString();

        const parser = new ComposerJsonExtractor();

        const actual = parser.extract({ "composer.json": new ExtractorFile(content) });

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

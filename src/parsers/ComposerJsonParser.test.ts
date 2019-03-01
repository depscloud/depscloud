import {readFile} from "fs";
import {promisify} from "util";
import ComposerJsonParser from "./ComposerJsonParser";

const readFileAsync = promisify(readFile);

describe("ComposerJsonParser", () => {
    test("fullParse", async () => {
        const composerPath = require.resolve("./testdata/composer.json");
        const buffer = await readFileAsync(composerPath);
        const content = buffer.toString();

        const parser = new ComposerJsonParser();

        const actual = parser.parse(composerPath, content);

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

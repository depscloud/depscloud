import {readFile} from "fs";
import {promisify} from "util";
import ExtractorFile from "./ExtractorFile";
import BowerJsonExtractor from "./BowerJsonExtractor";

const readFileAsync = promisify(readFile);

describe("BowerJsonExtractor", () => {
    test("fullParse", async () => {
        const jsonPath = require.resolve("./testdata/bower.json");
        const buffer = await readFileAsync(jsonPath);
        const content = buffer.toString();

        const parser = new BowerJsonExtractor();

        const actual = await parser.extract("", { "bower.json": new ExtractorFile(content) });

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

import {readFile} from "fs";
import {promisify} from "util";
import BowerJsonExtractor from "./BowerJsonExtractor";
import ExtractorFile from "./ExtractorFile";

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

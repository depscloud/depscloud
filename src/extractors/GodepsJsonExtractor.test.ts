import {readFile} from "fs";
import {promisify} from "util";
import ExtractorFile from "./ExtractorFile";
import GodepsJsonExtractor from "./GodepsJsonExtractor";

const readFileAsync = promisify(readFile);

describe("GodepsJsonExtractor", () => {
    test("fullParse", async () => {
        const godepsPath = require.resolve("./testdata/Godeps.json");
        const buffer = await readFileAsync(godepsPath);
        const content = buffer.toString();

        const parser = new GodepsJsonExtractor();

        const actual = parser.extract({ "Godeps.json": new ExtractorFile(content) });

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

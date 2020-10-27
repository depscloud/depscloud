import {readFile} from "fs";
import {promisify} from "util";
import ExtractorFile from "./ExtractorFile";
import JsonnetJsonExtractor from "./JsonnetJsonExtractor";

const readFileAsync = promisify(readFile);

describe("JsonnetJsonExtractor", () => {
    test("fullParse", async () => {
        const jsonPath = require.resolve("./testdata/jsonnetfile.json");
        const buffer = await readFileAsync(jsonPath);
        const content = buffer.toString();

        const parser = new JsonnetJsonExtractor();

        const actual = await parser.extract("", { "jsonnetfile.json": new ExtractorFile(content) });

        expect(actual).toMatchSnapshot();
    });
});

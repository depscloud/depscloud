import {readFile} from "fs";
import {promisify} from "util";
import ExtractorFile from "./ExtractorFile";
import JsonnetfileJsonExtractor from "./JsonnetfileJsonExtractor";

const readFileAsync = promisify(readFile);

describe("JsonnetfileJsonExtractor", () => {
    test("fullParse", async () => {
        const jsonPath = require.resolve("./testdata/jsonnetfile.json");
        const buffer = await readFileAsync(jsonPath);
        const content = buffer.toString();

        const parser = new JsonnetfileJsonExtractor();

        const actual = await parser.extract("", { "jsonnetfile.json": new ExtractorFile(content) });

        expect(actual).toMatchSnapshot();
    });
});

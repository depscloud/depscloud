import {readFile} from "fs";
import {promisify} from "util";
import ExtractorFile from "./ExtractorFile";
import PackageJsonExtractor from "./PackageJsonExtractor";

const readFileAsync = promisify(readFile);

describe("PackageJsonExtractor", () => {
    test("fullParse", async () => {
        const jsonPath = require.resolve("./testdata/package.json");
        const buffer = await readFileAsync(jsonPath);
        const content = buffer.toString();

        const parser = new PackageJsonExtractor();

        const actual = await parser.extract("", { "package.json": new ExtractorFile(content) });

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

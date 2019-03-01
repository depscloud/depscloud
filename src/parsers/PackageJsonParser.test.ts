import {readFile} from "fs";
import {promisify} from "util";
import PackageJsonParser from "./PackageJsonParser";

const readFileAsync = promisify(readFile);

describe("PackageJsonParser", () => {
    test("fullParse", async () => {
        const jsonPath = require.resolve("./testdata/package.json");
        const buffer = await readFileAsync(jsonPath);
        const content = buffer.toString();

        const parser = new PackageJsonParser();

        const actual = parser.parse(jsonPath, content);

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

import {readFile} from "fs";
import {promisify} from "util";
import GoModParser from "./GoModParser";

const readFileAsync = promisify(readFile);

describe("GoModParser", () => {
    test("fullParse", async () => {
        const godepsPath = require.resolve("./testdata/go.mod");
        const buffer = await readFileAsync(godepsPath);
        const content = buffer.toString();

        const parser = new GoModParser();

        const actual = parser.parse(godepsPath, content);

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

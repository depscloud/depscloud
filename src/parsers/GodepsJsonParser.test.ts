import {readFile} from "fs";
import {promisify} from "util";
import GodepsJsonParser from "./GodepsJsonParser";

const readFileAsync = promisify(readFile);

describe("GodepsJsonParser", () => {
    test("fullParse", async () => {
        const godepsPath = require.resolve("./testdata/Godeps.json");
        const buffer = await readFileAsync(godepsPath);
        const content = buffer.toString();

        const parser = new GodepsJsonParser();

        const actual = parser.parse(godepsPath, content);

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

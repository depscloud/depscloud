import {readFile} from "fs";
import {promisify} from "util";
import GopkgTomlParser from "./GopkgTomlParser";

const readFileAsync = promisify(readFile);

describe("GopkgTomlParser", () => {
    test("fullParse", async () => {
        const gopkgToml = require.resolve("./testdata/Gopkg.toml");
        const buffer = await readFileAsync(gopkgToml);
        const content = buffer.toString();

        const parser = new GopkgTomlParser();

        const actual = parser.parse(gopkgToml, content);

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

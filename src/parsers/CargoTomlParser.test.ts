import {readFile} from "fs";
import {promisify} from "util";
import CargoTomlParser from "./CargoTomlParser";

const readFileAsync = promisify(readFile);

describe("CargoTomlParser", () => {
    test("fullParse", async () => {
        const gopkgToml = require.resolve("./testdata/Cargo.toml");
        const buffer = await readFileAsync(gopkgToml);
        const content = buffer.toString();

        const parser = new CargoTomlParser();

        const actual = parser.parse(gopkgToml, content);

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

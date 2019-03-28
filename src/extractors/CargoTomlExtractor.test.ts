import {readFile} from "fs";
import {promisify} from "util";
import CargoTomlExtractor from "./CargoTomlExtractor";
import ExtractorFile from "./ExtractorFile";

const readFileAsync = promisify(readFile);

describe("CargoTomlExtractor", () => {
    test("fullParse", async () => {
        const gopkgToml = require.resolve("./testdata/Cargo.toml");
        const buffer = await readFileAsync(gopkgToml);
        const content = buffer.toString();

        const parser = new CargoTomlExtractor();

        const actual = parser.extract({ "Cargo.toml": new ExtractorFile(content) });

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

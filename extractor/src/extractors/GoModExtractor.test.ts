import {readFile} from "fs";
import {promisify} from "util";
import ExtractorFile from "./ExtractorFile";
import GoModExtractor from "./GoModExtractor";

const readFileAsync = promisify(readFile);

describe("GoModExtractor", () => {
    test("fullParse", async () => {
        const gomodPath = require.resolve("./testdata/go.mod");
        const buffer = await readFileAsync(gomodPath);
        const content = buffer.toString();

        const parser = new GoModExtractor();

        const actual = await parser.extract("", { "go.mod": new ExtractorFile(content) });

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

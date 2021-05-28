import {readFile} from "fs";
import {promisify} from "util";
import ExtractorFile from "./ExtractorFile";
import PipfileExtractor from "./PipfileExtractor";

const readFileAsync = promisify(readFile);

describe("PipfileExtractor", () => {
    test("fullParse", async () => {
        const filePath = require.resolve("./testdata/Pipfile");
        const buffer = await readFileAsync(filePath);
        const content = buffer.toString();

        const parser = new PipfileExtractor();

        const actual = await parser.extract("", { "Pipfile": new ExtractorFile(content) });

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

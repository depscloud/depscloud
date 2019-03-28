import {readFile} from "fs";
import {promisify} from "util";
import ExtractorFile from "./ExtractorFile";
import IvyXmlExtractor from "./IvyXmlExtractor";

const readFileAsync = promisify(readFile);

describe("IvyXmlExtractor", () => {
    test("fullParse", async () => {
        const ivyPath = require.resolve("./testdata/ivy.xml");
        const buffer = await readFileAsync(ivyPath);
        const content = buffer.toString();

        const parser = new IvyXmlExtractor();

        const actual = parser.extract({ "ivy.xml": new ExtractorFile(content) });

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

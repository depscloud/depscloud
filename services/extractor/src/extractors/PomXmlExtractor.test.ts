import {readFile} from "fs";
import {promisify} from "util";
import ExtractorFile from "./ExtractorFile";
import PomXmlExtractor from "./PomXmlExtractor";

const readFileAsync = promisify(readFile);

describe("PomXmlExtractor", () => {
    test("fullParse", async () => {
        const pomPath = require.resolve("./testdata/pom.xml");
        const buffer = await readFileAsync(pomPath);
        const content = buffer.toString();

        const parser = new PomXmlExtractor();

        const actual = await parser.extract("", { "pom.xml": new ExtractorFile(content) });

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

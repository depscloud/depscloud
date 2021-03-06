import {readFile} from "fs";
import {promisify} from "util";
import ExtractorFile from "./ExtractorFile";
import RequirementsTxtExtractor from "./RequirementsTxtExtractor";

const readFileAsync = promisify(readFile);

describe("RequirementsTxtExtractor", () => {
    test("fullParse", async () => {
        const filePath = require.resolve("./testdata/requirements.txt");
        const buffer = await readFileAsync(filePath);
        const content = buffer.toString();

        const parser = new RequirementsTxtExtractor();

        const url = "git@example.com:common/python-web.git";
        const actual = await parser.extract(url, { "requirements.txt": new ExtractorFile(content) });

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

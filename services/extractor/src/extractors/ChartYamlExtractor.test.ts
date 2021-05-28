import {readFile} from "fs";
import {promisify} from "util";
import ChartYamlExtractor from "./ChartYamlExtractor";
import ExtractorFile from "./ExtractorFile";

const readFileAsync = promisify(readFile);

describe("ChartYamlExtractor", () => {
    test("fullParse", async () => {
        const filePath = require.resolve("./testdata/Chart.yaml");
        const buffer = await readFileAsync(filePath);
        const content = buffer.toString();

        const parser = new ChartYamlExtractor();

        const actual = await parser.extract("", { "Chart.yaml": new ExtractorFile(content) });

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

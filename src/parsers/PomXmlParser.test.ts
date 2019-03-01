import {readFile} from "fs";
import {promisify} from "util";
import PomXmlParser from "./PomXmlParser";

const readFileAsync = promisify(readFile);

describe("PomXmlParser", () => {
    test("fullParse", async () => {
        const pomPath = require.resolve("./testdata/pom.xml");
        const buffer = await readFileAsync(pomPath);
        const content = buffer.toString();

        const parser = new PomXmlParser();

        const actual = parser.parse(pomPath, content);

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

import {readFile} from "fs";
import {promisify} from "util";
import IvyXmlParser from "./IvyXmlParser";

const readFileAsync = promisify(readFile);

describe("IvyXmlParser", () => {
    test("fullParse", async () => {
        const ivyPath = require.resolve("./testdata/ivy.xml");
        const buffer = await readFileAsync(ivyPath);
        const content = buffer.toString();

        const parser = new IvyXmlParser();

        const actual = parser.parse(ivyPath, content);

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

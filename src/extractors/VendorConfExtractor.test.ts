import {readFile} from "fs";
import {promisify} from "util";
import ExtractorFile from "./ExtractorFile";
import VendorConfExtractor from "./VendorConfExtractor";

const readFileAsync = promisify(readFile);

describe("VendorConfExtractor", () => {
    test("fullParse", async () => {
        const vendorConfPath = require.resolve("./testdata/vendor.conf");
        const buffer = await readFileAsync(vendorConfPath);
        const content = buffer.toString();

        const parser = new VendorConfExtractor();

        const actual = await parser.extract("", { "vendor.conf": new ExtractorFile(content) });

        expect(actual).toMatchSnapshot();
        expect(JSON.stringify(actual, null, 2)).toMatchSnapshot();
    });
});

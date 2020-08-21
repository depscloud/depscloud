import cheerio = require("cheerio");
import YAML = require("js-yaml");
import TOML = require("toml");

export default class ExtractorFile {
    private readonly body: string;

    constructor(body: string) {
        this.body = body;
    }

    public raw(): string {
        return this.body;
    }

    public json(): any {
        return JSON.parse(this.body);
    }

    public toml(): any {
        return TOML.parse(this.body);
    }

    public yaml(): any {
        return YAML.safeLoad(this.body);
    }

    public xml(): Cheerio {
        return cheerio.load(this.body).root();
    }
}

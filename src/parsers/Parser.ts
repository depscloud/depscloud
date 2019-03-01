import cheerio = require("cheerio");
import {DependencyManagementFile} from "../../api/deps";

export interface IParser {
    pathMatch(path: string): boolean;
    parse(path: string, content: string): DependencyManagementFile;
}

export class CompositeParser implements IParser {
    private readonly parsers: IParser[];

    constructor(parsers: IParser[]) {
        this.parsers = parsers;
    }

    public pathMatch(path: string): boolean {
        return this.parsers.some((parser) => parser.pathMatch(path));
    }

    public parse(path: string, content: string): DependencyManagementFile {
        for (const parser of this.parsers) {
            if (parser.pathMatch(path)) {
                return parser.parse(path, content);
            }
        }
        return null;
    }
}

export class XmlParser implements IParser {
    public parse(path: string, content: string): DependencyManagementFile {
        return this.parseXml(cheerio.load(content).root());
    }

    public pathMatch(path: string): boolean {
        throw new MethodNotImplementedError("pathMatch");
    }

    public parseXml(xml: Cheerio): DependencyManagementFile {
        throw new MethodNotImplementedError("parseXml");
    }
}

export class JsonParser implements IParser {
    public parse(path: string, content: string): DependencyManagementFile {
        return this.parseJson(JSON.parse(content));
    }

    public pathMatch(path: string): boolean {
        throw new MethodNotImplementedError("pathMatch");
    }

    public parseJson(json: any): DependencyManagementFile {
        throw new MethodNotImplementedError("parseJson");
    }

}

class MethodNotImplementedError extends Error {
    constructor(methodName: string) {
        super(`MethodNotImplemented: You must override the default behavior for ${methodName}`);
    }
}

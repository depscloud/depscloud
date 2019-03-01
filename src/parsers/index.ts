import ComposerJsonParser from "./ComposerJsonParser";
import GodepsJsonParser from "./GodepsJsonParser";
import IvyXmlParser from "./IvyXmlParser";
import PackageJsonParser from "./PackageJsonParser";
import {CompositeParser, IParser} from "./Parser";
import PomXmlParser from "./PomXmlParser";

export function defaultParser(): IParser {
    return new CompositeParser([
        new ComposerJsonParser(),
        new GodepsJsonParser(),
        new IvyXmlParser(),
        new PackageJsonParser(),
        new PomXmlParser(),
    ]);
}

import CargoTomlParser from "./CargoTomlParser";
import ComposerJsonParser from "./ComposerJsonParser";
import GodepsJsonParser from "./GodepsJsonParser";
import GoModParser from "./GoModParser";
import IvyXmlParser from "./IvyXmlParser";
import PackageJsonParser from "./PackageJsonParser";
import {CompositeParser, IParser} from "./Parser";
import PomXmlParser from "./PomXmlParser";

export function defaultParser(): IParser {
    return new CompositeParser([
        new CargoTomlParser(),
        new ComposerJsonParser(),
        new GodepsJsonParser(),
        new GoModParser(),
        new IvyXmlParser(),
        new PackageJsonParser(),
        new PomXmlParser(),
    ]);
}

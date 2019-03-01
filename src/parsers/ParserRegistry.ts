import Registry from "../common/Registry";
import ComposerJsonParser from "./ComposerJsonParser";
import GodepsJsonParser from "./GodepsJsonParser";
import GoModParser from "./GoModParser";
import IvyXmlParser from "./IvyXmlParser";
import PackageJsonParser from "./PackageJsonParser";
import {CompositeParser, IParser} from "./Parser";
import PomXmlParser from "./PomXmlParser";

const ParserRegistry = new Registry<IParser>("Parser");

ParserRegistry.registerAll({
    CompositeParser: async (params) => new CompositeParser(params),
    ComposerJsonParser: async (_) => new ComposerJsonParser(),
    GodepsJsonParser: async (_) => new GodepsJsonParser(),
    GoModParser: async (_) => new GoModParser(),
    IvyXmlParser: async (_) => new IvyXmlParser(),
    PackageJsonParser: async (_) => new PackageJsonParser(),
    PomXmlParser: async (_) => new PomXmlParser(),
});

export default ParserRegistry;

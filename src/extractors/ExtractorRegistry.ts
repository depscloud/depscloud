import Registry from "../common/Registry";
import CargoTomlExtractor from "./CargoTomlExtractor";
import ComposerJsonExtractor from "./ComposerJsonExtractor";
import Extractor from "./Extractor";
import GodepsJsonExtractor from "./GodepsJsonExtractor";
import GoModExtractor from "./GoModExtractor";
import GopkgTomlExtractor from "./GopkgTomlExtractor";
import IvyXmlExtractor from "./IvyXmlExtractor";
import PackageJsonExtractor from "./PackageJsonExtractor";
import PomXmlExtractor from "./PomXmlExtractor";

const ExtractorRegistry = new Registry<Extractor>("Extractor");

ExtractorRegistry.registerAll({
    CargoTomlParser: async (_) => new CargoTomlExtractor(),
    ComposerJsonParser: async (_) => new ComposerJsonExtractor(),
    GodepsJsonParser: async (_) => new GodepsJsonExtractor(),
    GoModParser: async (_) => new GoModExtractor(),
    GopkgTomlParser: async (_) => new GopkgTomlExtractor(),
    IvyXmlParser: async (_) => new IvyXmlExtractor(),
    PackageJsonParser: async (_) => new PackageJsonExtractor(),
    PomXmlParser: async (_) => new PomXmlExtractor(),
});

export default ExtractorRegistry;

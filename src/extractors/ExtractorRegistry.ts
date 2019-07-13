import Registry from "../common/Registry";
import BuildGradleExtractor from "./BuildGradleExtractor";
import CargoTomlExtractor from "./CargoTomlExtractor";
import ComposerJsonExtractor from "./ComposerJsonExtractor";
import Extractor from "./Extractor";
import GodepsJsonExtractor from "./GodepsJsonExtractor";
import GoModExtractor from "./GoModExtractor";
import GopkgTomlExtractor from "./GopkgTomlExtractor";
import IvyXmlExtractor from "./IvyXmlExtractor";
import PackageJsonExtractor from "./PackageJsonExtractor";
import PomXmlExtractor from "./PomXmlExtractor";
import BowerJsonExtractor from "./BowerJsonExtractor";

const ExtractorRegistry = new Registry<Extractor>("Extractor");

ExtractorRegistry.registerAll({
    BuildGradleExtractor: async (_) => new BuildGradleExtractor(),
    CargoTomlExtractor: async (_) => new CargoTomlExtractor(),
    ComposerJsonExtractor: async (_) => new ComposerJsonExtractor(),
    GodepsJsonExtractor: async (_) => new GodepsJsonExtractor(),
    GoModExtractor: async (_) => new GoModExtractor(),
    GopkgTomlExtractor: async (_) => new GopkgTomlExtractor(),
    IvyXmlExtractor: async (_) => new IvyXmlExtractor(),
    PackageJsonExtractor: async (_) => new PackageJsonExtractor(),
    PomXmlExtractor: async (_) => new PomXmlExtractor(),
    BowerJsonExtractor: async (_) => new BowerJsonExtractor(),
});

export default ExtractorRegistry;

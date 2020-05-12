import Registry from "../common/Registry";
import BowerJsonExtractor from "./BowerJsonExtractor";
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
import VendorConfExtractor from "./VendorConfExtractor";

const ExtractorRegistry = new Registry<Extractor>("Extractor");

ExtractorRegistry.registerAll({
    "build.gradle": async (_) => new BuildGradleExtractor(),
    "Cargo.toml": async (_) => new CargoTomlExtractor(),
    "composer.json": async (_) => new ComposerJsonExtractor(),
    "Godeps.json": async (_) => new GodepsJsonExtractor(),
    "go.mod": async (_) => new GoModExtractor(),
    "Gopkg.toml": async (_) => new GopkgTomlExtractor(),
    "ivy.xml": async (_) => new IvyXmlExtractor(),
    "package.json": async (_) => new PackageJsonExtractor(),
    "pom.xml": async (_) => new PomXmlExtractor(),
    "bower.json": async (_) => new BowerJsonExtractor(),
    "vendor.conf": async (_) => new VendorConfExtractor(),
});

export default ExtractorRegistry;

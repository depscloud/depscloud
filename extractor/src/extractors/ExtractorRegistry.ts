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
import JsonnetfileJsonExtractor from "./JsonnetfileJsonExtractor";
import PipfileExtractor from "./PipfileExtractor";
import RequirementsTxtExtractor from "./RequirementsTxtExtractor";

const ExtractorRegistry = new Registry<Extractor>("Extractor");

ExtractorRegistry.registerAll({
    "build.gradle": async () => new BuildGradleExtractor(),
    "Cargo.toml": async () => new CargoTomlExtractor(),
    "composer.json": async () => new ComposerJsonExtractor(),
    "Godeps.json": async () => new GodepsJsonExtractor(),
    "go.mod": async () => new GoModExtractor(),
    "Gopkg.toml": async () => new GopkgTomlExtractor(),
    "ivy.xml": async () => new IvyXmlExtractor(),
    "package.json": async () => new PackageJsonExtractor(),
    "pom.xml": async () => new PomXmlExtractor(),
    "bower.json": async () => new BowerJsonExtractor(),
    "vendor.conf": async () => new VendorConfExtractor(),
    "jsonnetfile.json": async() => new JsonnetfileJsonExtractor(),
    "Pipfile": async() => new PipfileExtractor(),
    "requirements.txt": async() => new RequirementsTxtExtractor(),
});

export default ExtractorRegistry;

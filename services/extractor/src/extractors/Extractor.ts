import {ManifestFile} from "@depscloud/api/v1beta";
import ExtractorFile from "./ExtractorFile";
import MatchConfig from "../matcher/MatchConfig";

export default interface Extractor {
    matchConfig(): MatchConfig;
    requires(): string[];
    extract(url: string, files: { [key: string]: ExtractorFile }): Promise<ManifestFile>;
}

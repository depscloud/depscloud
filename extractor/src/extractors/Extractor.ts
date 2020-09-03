import {DependencyManagementFile} from "@depscloud/api/v1alpha/deps";
import ExtractorFile from "./ExtractorFile";
import MatchConfig from "../matcher/MatchConfig";

export default interface Extractor {
    matchConfig(): MatchConfig;
    requires(): string[];
    extract(url: string, files: { [key: string]: ExtractorFile }): Promise<DependencyManagementFile>;
}

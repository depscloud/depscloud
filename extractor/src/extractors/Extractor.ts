import {DependencyManagementFile} from "@depscloud/api/v1alpha/deps";
import ExtractorFile from "./ExtractorFile";

export default interface Extractor {
    requires(): string[];
    extract(url: string, files: { [key: string]: ExtractorFile }): Promise<DependencyManagementFile>;
}

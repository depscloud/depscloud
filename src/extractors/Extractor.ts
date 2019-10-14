import {DependencyManagementFile} from "@deps-cloud/api/v1alpha/deps/deps";
import ExtractorFile from "./ExtractorFile";

export default interface Extractor {
    requires(): string[];
    extract(url: string, files: { [key: string]: ExtractorFile }): Promise<DependencyManagementFile>;
}

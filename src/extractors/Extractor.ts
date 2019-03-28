import {DependencyManagementFile} from "../../api/deps";
import ExtractorFile from "./ExtractorFile";

export default interface Extractor {
    requires(): string[];
    extract(files: { [key: string]: ExtractorFile }): DependencyManagementFile;
}

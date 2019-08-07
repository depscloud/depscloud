import {ServerUnaryCall} from "grpc";
import {DependencyManagementFile} from "../../api/deps";
import {ExtractRequest, ExtractResponse, MatchRequest, MatchResponse} from "../../api/extractor";
import Extractor from "../extractors/Extractor";
import ExtractorFile from "../extractors/ExtractorFile";
import AsyncDependencyExtractor from "./AsyncDependencyExtractor";

function constructTree(separator: string, paths: string[]): any {
    const root: any = {};

    paths.forEach((key) => {
        const parts = key.split(separator);

        let ptr = root;
        let i = 0;
        for (; i < parts.length - 1; i++) {
            const part = parts[i];
            if (!ptr[part]) {
                ptr[part] = {};
            }
            ptr = ptr[part];
        }
        ptr[parts[i]] = key;
    });

    return root;
}

export default class DependencyExtractorImpl implements AsyncDependencyExtractor {
    private readonly extractors: Extractor[];

    constructor(extractors: Extractor[]) {
        this.extractors = extractors;
    }

    public matchInternal(separator: string, paths: string[]): string[] {
        const root = constructTree(separator, paths);

        let level = [ root ];
        const matchedPaths = [];

        while (level.length > 0) {
            const size = level.length;

            for (let i = 0; i < size; i++) {
                const dir = level.shift();

                this.extractors
                    .filter((extractor) =>
                        extractor.requires()
                            .map((req) => dir[req] && typeof dir[req] === "string")
                            .reduce((last, current) => last && current))
                    .forEach((extractor) =>
                        extractor.requires()
                            .forEach((req) => matchedPaths.push(dir[req])));

                const nextLevel = Object.keys(dir)
                    .map((name) => dir[name])
                    .filter((val) => typeof val !== "string");

                level = level.concat(nextLevel);
            }
        }

        return matchedPaths;
    }

    public async match(call: ServerUnaryCall<MatchRequest>): Promise<MatchResponse> {
        const { separator, paths } = call.request;

        return {
            matchedPaths: this.matchInternal(separator, paths),
        };
    }

    public async extractInternal(
        url: string,
        separator: string,
        fileContents: { [key: string]: string },
    ): Promise<DependencyManagementFile[]> {
        const matchedPaths = this.matchInternal(separator, Object.keys(fileContents));

        const root = constructTree(separator, matchedPaths);

        let level = [ root ];
        let managementFilePromises: Array<Promise<DependencyManagementFile>> = [];

        while (level.length > 0) {
            const size = level.length;

            for (let i = 0; i < size; i++) {
                const dir = level.shift();

                const nextManagementFilePromises = this.extractors
                    .filter((extractor) =>
                        extractor.requires()
                            .map((req) => dir[req] && typeof dir[req] === "string")
                            .reduce((last, current) => last && current))
                    .map((extractor) => {
                        const files = {};
                        extractor.requires().forEach((req) => {
                            const key = dir[req];
                            files[req] = new ExtractorFile(fileContents[key]);
                        });
                        return extractor.extract(url, files);
                    });

                managementFilePromises = managementFilePromises.concat(nextManagementFilePromises);

                const nextLevel = Object.keys(dir)
                    .map((name) => dir[name])
                    .filter((val) => typeof val !== "string");

                level = level.concat(nextLevel);
            }
        }

        let managementFiles = await Promise.all(managementFilePromises);
        managementFiles = managementFiles.filter((f) => !!f);

        return managementFiles;
    }

    public async extract(call: ServerUnaryCall<ExtractRequest>): Promise<ExtractResponse> {
        const { url, separator, fileContents } = call.request;

        const managementFiles = await this.extractInternal(url, separator, fileContents);

        return {
            managementFiles,
        };
    }
}

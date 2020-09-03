import {DependencyManagementFile} from "@depscloud/api/v1alpha/deps";
import {
    ExtractRequest, ExtractResponse, MatchRequest, MatchResponse,
} from "@depscloud/api/v1alpha/extractor";
import {ServerUnaryCall} from "@grpc/grpc-js";
import ExtractorFile from "../extractors/ExtractorFile";
import AsyncDependencyExtractor from "./AsyncDependencyExtractor";
import MatcherAndExtractor from "./MatcherAndExtractor";

import path = require("path")

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

function normalizePaths(separator: string, paths: string[]): string[] {
    return paths.map((p) => {
        if (separator === path.win32.sep) {
            return path.posix.join(path.win32.parse(p))
        }
        return p
    })
}

export default class DependencyExtractorImpl implements AsyncDependencyExtractor {
    private readonly matcherAndExtractors: MatcherAndExtractor[];

    constructor(matcherAndExtractors: MatcherAndExtractor[]) {
        this.matcherAndExtractors = matcherAndExtractors;
    }

    public matchInternal(separator: string, paths: string[]): string[] {
        const matchedPaths = [];
        const normPaths = normalizePaths(separator, paths);

        normPaths.forEach((p, i) => {
            const found = this.matcherAndExtractors.find((me) => me.matcher.match(p))
            if (found) {
                matchedPaths.push(paths[i]);
            }
        })

        return matchedPaths;
    }

    public async match(call: ServerUnaryCall<MatchRequest, MatchResponse>): Promise<MatchResponse> {
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
        const paths = Object.keys(fileContents);
        const matchedPaths = this.matchInternal(separator, paths);

        const root = constructTree(separator, matchedPaths);

        let level = [ root ];
        let managementFilePromises: Promise<DependencyManagementFile>[] = [];

        while (level.length > 0) {
            const size = level.length;

            for (let i = 0; i < size; i++) {
                const dir = level.shift();

                const nextManagementFilePromises = this.matcherAndExtractors
                    .filter((me) =>
                        me.extractor.requires()
                            .map((req) => dir[req] && typeof dir[req] === "string")
                            .reduce((last, current) => last && current))
                    .map((me) => {
                        const files = {};
                        me.extractor.requires().forEach((req) => {
                            const key = dir[req];
                            files[req] = new ExtractorFile(fileContents[key]);
                        });
                        return me.extractor.extract(url, files);
                    });

                managementFilePromises = managementFilePromises.concat(nextManagementFilePromises);

                const nextLevel = Object.keys(dir)
                    .map((name) => dir[name])
                    .filter((val) => typeof val !== "string");

                level = level.concat(nextLevel);
            }
        }

        let managementFiles = await Promise.all(managementFilePromises);
        managementFiles = managementFiles
            .filter((f) => !!f)         // ensure no nulls returned
            .filter((f) => !!f.module); // ensure a module is returned

        return managementFiles;
    }

    public async extract(call: ServerUnaryCall<ExtractRequest, ExtractResponse>): Promise<ExtractResponse> {
        const { url, separator, fileContents } = call.request;

        const managementFiles = await this.extractInternal(url, separator, fileContents);

        return {
            managementFiles,
        };
    }
}

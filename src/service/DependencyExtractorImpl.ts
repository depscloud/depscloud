import fs = require("fs");
import {ServerUnaryCall} from "grpc";
import { Clone, Cred } from "nodegit";
import path = require("path");
import tmp = require("tmp");
import { DependencyManagementFile } from "../../api/deps";
import { ExtractRequest, ExtractResponse } from "../../api/extractor";
import {IParser} from "../parsers/Parser";
import AsyncDependencyExtractor from "./AsyncDependencyExtractor";
const fsp = fs.promises;

function tmpdir(): Promise<[string, () => void]> {
    return new Promise((resolve, reject) => {
        tmp.dir({
            unsafeCleanup: true,
        }, (err, tmpPath, cleanup) => {
            if (err != null) {
                reject(err);
            } else {
                resolve([ tmpPath, cleanup ]);
            }
        });
    });
}

export default class DependencyExtractorImpl implements AsyncDependencyExtractor {
    private readonly parser: IParser;
    private readonly credentials: (url: string, username: string) => Promise<Cred>;

    constructor(parser: IParser, credentials: (url: string, username: string) => Promise<Cred>) {
        this.parser = parser;
        this.credentials = credentials;
    }

    private async read({ path: filePath }: any, collector: DependencyManagementFile[]): Promise<void> {
        const contents = await fsp.readFile(filePath);
        const dependencyManagementFile = this.parser.parse(filePath, contents.toString());
        collector.push(dependencyManagementFile);
    }

    private async walk(url: string, collector: DependencyManagementFile[]): Promise<void> {
        const files = await fsp.readdir(url);
        const fileStatPromises: Array<Promise<any>> = files.map((file: string) => {
            return {
                name: file,
                path: path.join(url, file),
            };
        }).map((file: any) => {
            return fsp.stat(file.path)
                .then((stats) => {
                    file.stats = stats;
                    return file;
                });
        });

        const fileStats = await Promise.all(fileStatPromises);

        const outstanding = fileStats
            .filter((fileStat) => fileStat != null)
            .map((fileStat) => {
                if (fileStat.stats.isDirectory()) {
                    return this.walk(fileStat.path, collector);
                } else if (this.parser.pathMatch(fileStat.path)) {
                    return this.read(fileStat, collector);
                }
            });

        await Promise.all(outstanding);
    }

    public async extract(call: ServerUnaryCall<ExtractRequest>): Promise<ExtractResponse> {
        const managementFiles: DependencyManagementFile[] = [];
        const [ tmpPath, cleanup ] = await tmpdir();
        const repositoryUrl = call.request.repositoryUrl;

        try {
            await Clone.clone(repositoryUrl, tmpPath, {
                fetchOpts: {
                    callbacks: {
                        certificateCheck: () => 1,
                        credentials: this.credentials,
                    },
                },
            });

            await this.walk(tmpPath, managementFiles);
        } finally {
            cleanup();
        }

        return {
            repositoryUrl,
            managementFiles,
        };
    }
}

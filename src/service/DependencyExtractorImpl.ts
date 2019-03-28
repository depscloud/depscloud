import {ServerUnaryCall} from "grpc";
import {Clone, Cred} from "nodegit";
import {DependencyManagementFile} from "../../api/deps";
import {ExtractRequest, ExtractResponse} from "../../api/extractor";
import Extractor from "../extractors/Extractor";
import ExtractorFile from "../extractors/ExtractorFile";
import AsyncDependencyExtractor from "./AsyncDependencyExtractor";

import fs = require("fs");
import path = require("path");
import tmp = require("tmp");

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
    private readonly extractors: Extractor[];
    private readonly credentials: (url: string, username: string) => Promise<Cred>;

    constructor(extractors: Extractor[], credentials: (url: string, username: string) => Promise<Cred>) {
        this.extractors = extractors;
        this.credentials = credentials;
    }

    private async checkDirectory(
        dir: { [key: string]: [ string, fs.Stats ] },
        extractor: Extractor,
    ): Promise<DependencyManagementFile> {
        const meetsRequirements = extractor.requires()
            .map((req) => dir[req][1].isFile())
            .reduce((last, current) => last && current);

        if (!meetsRequirements) {
            return null;
        }

        const files: { [key: string]: ExtractorFile } = {};

        const readFileReqs = extractor.requires()
            .map((req) => {
                return fsp.readFile(dir[req][0])
                    .then((buf) => {
                        files[req] = new ExtractorFile(buf.toString());
                    });
            });

        await Promise.all(readFileReqs);

        return extractor.extract(files);
    }

    private async walk(url: string): Promise<DependencyManagementFile[]> {
        const dir: { [key: string]: [ string, fs.Stats ] } = {};

        const fileReqs = await fsp.readdir(url)
            .then((readdir) => {
                return readdir.map((fileName) => {
                    const filePath = path.join(url, fileName);

                    return fsp.stat(filePath)
                        .then((fileStat) => {
                            dir[fileName] = [ filePath, fileStat ];
                        });
                });
            });

        await Promise.all(fileReqs);

        const promises: Array<Promise<DependencyManagementFile>> = [];
        for (const extractor of this.extractors) {
            promises.push(this.checkDirectory(dir, extractor));
        }

        const walkReqs: Array<Promise<DependencyManagementFile[]>> = Object.keys(dir)
            .map((fileName) => ({
                fileName,
                filePath: dir[fileName][0],
                fileStat: dir[fileName][1],
            }))
            .filter(({ fileStat }) => fileStat.isDirectory())
            .map(({ filePath }) => this.walk(filePath));

        walkReqs.push(Promise.all(promises));

        const walk: DependencyManagementFile[][] = await Promise.all(walkReqs);
        const deps: DependencyManagementFile[] = [].concat(...walk);

        return deps.filter((e) => e != null);
    }

    public async extract(call: ServerUnaryCall<ExtractRequest>): Promise<ExtractResponse> {
        let managementFiles: DependencyManagementFile[] = [];
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

            managementFiles = await this.walk(tmpPath);
        } finally {
            cleanup();
        }

        return {
            repositoryUrl,
            managementFiles,
        };
    }
}

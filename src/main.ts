import {promises} from "fs";
import {Server, ServerCredentials} from "grpc";
import {configure, getLogger} from "log4js";
import {Cred} from "nodegit";
import {DependencyExtractor} from "../api/extractor";
import {defaultParser} from "./parsers";
import AsyncDependencyExtractor from "./service/AsyncDependencyExtractor";
import DependencyExtractorImpl from "./service/DependencyExtractorImpl";
import unasyncify from "./service/unasyncify";
import program = require("caporal");

const logger = getLogger();

program.name("finch-extractor")
    .option("--port <port>", "The port to bind to.", program.INT)
    .option("--public-key <file>", "The public key to use when cloning.", program.STRING)
    .option("--private-key <file>", "The private key to use when cloning.", program.STRING)
    .option("--passphrase <file>", "The passphrasee to the public and private key.", program.STRING)
    .action(async (args: any, options: any) => {
        configure({
            appenders: {
                console: { type: "console" },
            },
            categories: {
                default: {
                    appenders: [ "console" ],
                    level: "debug",
                },
            },
        });

        let passphrase = null;
        if (options.passphrase) {
            passphrase = await promises.readFile(options.passphrase).then((c) => c.toString())
        }

        const credentials = async (url: string, username: string): Promise<Cred> => {
            if (url.startsWith("git@")) {
                if (options.publicKey && options.privateKey) {
                    return Promise.resolve(Cred.sshKeyNew(username, options.publicKey, options.privateKey, passphrase));
                } else {
                    return Cred.sshKeyFromAgent(username);
                }
            }
            
            throw new Error(`unsupported url: ${url}`);
        };

        const port = options.port || 8090;
        const service: AsyncDependencyExtractor = new DependencyExtractorImpl(defaultParser(), credentials);

        const server = new Server();
        server.addService(DependencyExtractor.service, unasyncify(service));
        server.bind(`0.0.0.0:${port}`, ServerCredentials.createInsecure());
        logger.info(`[main] starting gRPC on :${port}`);
        server.start();
    })
    .parse(process.argv);

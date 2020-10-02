import {DependencyExtractor} from "@depscloud/api/v1alpha/extractor";

import {Server, ServerCredentials} from "@grpc/grpc-js";
import {configure, getLogger} from "log4js";
import ExtractorRegistry from "./extractors/ExtractorRegistry";
import AsyncDependencyExtractor from "./service/AsyncDependencyExtractor";
import DependencyExtractorImpl from "./service/DependencyExtractorImpl";
import unasyncify from "./service/unasyncify";

import express = require("express");
import program = require("caporal");
import fs = require("fs");
import health = require("grpc-health-check/health");
import healthv1 = require("grpc-health-check/v1/health_pb");
import Matcher from "./matcher/Matcher";
import promMiddleware = require('express-prometheus-middleware');

const asyncFs = fs.promises;

const logger = getLogger();

program.name("extractor")
    .option("--bind-address <bindAddress>", "The ip address to bind to.", program.STRING)
    .option("--http-port <port>", "The port to run http on.", program.INT)
    .option("--grpc-port <port>", "The port to bind to.", program.INT)
    .option("--port <port>", "The port to bind to.", program.INT)
    .option("--tls-key <key>", "The path to the private key used for TLS", program.STRING)
    .option("--tls-cert <cert>", "The path to the certificate used for TLS", program.STRING)
    .option("--tls-ca <ca>", "The path to the certificate authority used for TLS", program.STRING)
    .option("--disable-manifests <manifest>", "The manifests to disable support for", program.ARRAY)
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

        const disabledManifests = (options.disableManifests || [])
            .reduce((agg, item) => {
                agg[item] = true;
                return agg;
            }, {});

        const extractorReqs = ExtractorRegistry.known()
            .filter((e) => !disabledManifests[e])
            .map((extractor) => ExtractorRegistry.resolve(extractor, null));

        const extractors = await Promise.all(extractorReqs);

        const matchersAndExtractors = extractors.map((extractor) => {
            return {
                matcher: new Matcher(extractor.matchConfig()),
                extractor,
            }
        });

        const impl: AsyncDependencyExtractor = new DependencyExtractorImpl(matchersAndExtractors);

        const healthcheck = new health.Implementation({
            "": healthv1.HealthCheckResponse.ServingStatus.SERVING,
        });
        // toggle the service health as such
        // healthcheck.setStatus("", healthv1.HealthCheckResponse.ServingStatus.NOT_SERVING);

        const server = new Server();
        server.addService(DependencyExtractor.service, unasyncify(impl));
        server.addService(health.service, healthcheck);

        let credentials = ServerCredentials.createInsecure();
        if (options.tlsKey && options.tlsCert && options.tlsCa) {
            logger.info("[main] configuring tls");

            const [ key, cert, ca ] = await Promise.all([
                asyncFs.readFile(options.tlsKey),
                asyncFs.readFile(options.tlsCert),
                asyncFs.readFile(options.tlsCa),
            ]);

            credentials = ServerCredentials.createSsl(ca, [ {
                private_key: key,
                cert_chain: cert,
            } ], true);
        }

        const bindAddress = options.bindAddress || "0.0.0.0";
        const httpPort = options.httpPort || 8080;
        const grpcPort = options.port || options.grpcPort || 8090;

        server.bindAsync(`${bindAddress}:${grpcPort}`, credentials, (err) => {
            if (err != null) {
                logger.error(err);
                return;
            }

            logger.info(`[main] starting grpc on ${bindAddress}:${grpcPort}`);
            server.start();
        });

        const app = express();

        app.use(promMiddleware({
            metricsPath: "/metrics",
            collectDefaultMetrics: true,
        }))

        const healthHandle = (req, resp) => {
            resp.json({
                state: "ok",
                timestamp: new Date(),
                results: {},
            });
        };

        app.get("/healthz", healthHandle);
        app.get("/health", healthHandle);

        app.listen(httpPort, () => {
            logger.info(`[main] starting http on ${bindAddress}:${httpPort}`)
        });
    })
    .parse(process.argv);

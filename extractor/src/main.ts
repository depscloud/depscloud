import {DependencyExtractor} from "@depscloud/api/v1alpha/extractor";
import {ManifestExtractionService} from "@depscloud/api/v1beta";

import {Server, ServerCredentials} from "@grpc/grpc-js";
import {addLayout, configure, getLogger} from "log4js";
import ExtractorRegistry from "./extractors/ExtractorRegistry";
import AsyncManifestExtractionService from "./service/AsyncManifestExtractionService";
import ManifestExtractionServiceImpl from "./service/ManifestExtractionServiceImpl";
import unasyncify from "./service/unasyncify";

import express = require("express");
import program = require("caporal");
import fs = require("fs");
import health = require("grpc-health-check/health");
import healthv1 = require("grpc-health-check/v1/health_pb");
import Matcher from "./matcher/Matcher";
import promMiddleware = require("express-prometheus-middleware");

const packageMeta = require("../package.json");

const asyncFs = fs.promises;

addLayout("json", function() {
    return function(logEvent) {
        const data = logEvent.data.length > 1 ? logEvent.data[1] : {};

        return JSON.stringify({
            level: logEvent.level.levelStr.toLowerCase(),
            ts: new Date(logEvent.startTime).getTime(),
            msg: logEvent.data[0],
            ...data,
        });
    }
})

const logger = getLogger();

program.name("extractor")
    .option("--bind-address <bindAddress>", "the ip address to bind to", program.STRING)
    .option("--http-port <port>", "the port to run http on", program.INT)
    .option("--grpc-port <port>", "the port to bind to", program.INT)
    .option("--port <port>", "the port to bind to", program.INT)
    .option("--tls-key <key>", "the path to the private key used for TLS", program.STRING)
    .option("--tls-cert <cert>", "the path to the certificate used for TLS", program.STRING)
    .option("--tls-ca <ca>", "the path to the certificate authority used for TLS", program.STRING)
    .option("--disable-manifests <manifest>", "the manifests to disable support for", program.ARRAY)
    .option("--log-level <level>", "configures the level at with logs are written", program.STRING)
    .option("--log-format <format>", "configures the format of the logs (console / json)", program.STRING)
    .action(async (args: any, options: any) => {
        const logFormat = options.logFormat == "console" ? "basic" : "json";

        configure({
            appenders: {
                out: { type: "console", layout: { type: logFormat } },
            },
            categories: {
                default: {
                    appenders: [ "out" ],
                    level: options.logLevel || "info",
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

        const impl: AsyncManifestExtractionService = new ManifestExtractionServiceImpl(matchersAndExtractors);

        const healthcheck = new health.Implementation({
            "": healthv1.HealthCheckResponse.ServingStatus.SERVING,
        });
        // toggle the service health as such
        // healthcheck.setStatus("", healthv1.HealthCheckResponse.ServingStatus.NOT_SERVING);

        const server = new Server();
        server.addService(ManifestExtractionService.service, unasyncify(impl));
        server.addService(health.service, healthcheck);

        let credentials = ServerCredentials.createInsecure();
        if (options.tlsKey && options.tlsCert && options.tlsCa) {
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
                program.fatalError(err)
                return;
            }

            logger.info("starting server", {
                "protocol": "grpc",
                "bind": `${bindAddress}:${grpcPort}`,
                "tls": options.tlsKey && options.tlsCert && options.tlsCa,
            });
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

        app.get("/version", (req, resp) => {
		    resp.json(packageMeta.meta);
        });

        app.listen(httpPort, () => {
            logger.info("starting server", {
                "protocol": "http",
                "bind": `${bindAddress}:${httpPort}`,
                "tls": false,
            });
        });
    })
    .parse(process.argv);

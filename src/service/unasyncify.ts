import {IDependencyExtractor} from "@deps-cloud/api/v1alpha/extractor/extractor";
import {getLogger} from "log4js";
import AsyncDependencyExtractor from "./AsyncDependencyExtractor";

const logger = getLogger();

function toString(error: Error): string {
    if (error.stack) {
        return error.stack;
    }
    return `${error.name}: ${error.message}`;
}

export default function unasyncify(instance: AsyncDependencyExtractor): IDependencyExtractor {
    return {
        match(request, callback) {
            logger.trace(`[service] match request: ${JSON.stringify(request)}`);

            instance.match(request)
                .then((response) => callback(null, response))
                .catch((error) => {
                    logger.error(`[service] match error: ${toString(error)}`);
                    callback(error, null);
                });
        },
        extract(request, callback) {
            logger.trace(`[service] extract request: ${JSON.stringify(request)}`);

            instance.extract(request)
                .then((response) => callback(null, response))
                .catch((error) => {
                    logger.error(`[service] extract error: ${toString(error)}`);
                    callback(error, null);
                });
        },
    };
}

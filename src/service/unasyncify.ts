import {UntypedServiceImplementation} from "@grpc/grpc-js";
import {getLogger} from "log4js";
import AsyncDependencyExtractor from "./AsyncDependencyExtractor";

const logger = getLogger();

function toString(error: Error): string {
    if (error.stack) {
        return error.stack;
    }
    return `${error.name}: ${error.message}`;
}

export default function unasyncify(instance: AsyncDependencyExtractor): UntypedServiceImplementation {
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

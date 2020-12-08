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
        match(call, callback) {
            instance.match(call)
                .then((response) => callback(null, response))
                .catch((error) => {
                    logger.error("match error", {
                        err: toString(error)
                    });
                    callback(error, null);
                });
        },
        extract(call, callback) {
            instance.extract(call)
                .then((response) => callback(null, response))
                .catch((error) => {
                    logger.error("extract error", {
                        err: toString(error)
                    });
                    callback(error, null);
                });
        },
    };
}

import {getLogger} from "log4js";
import {IDependencyExtractor} from "../../api/extractor";
import AsyncDependencyExtractor from "./AsyncDependencyExtractor";

const logger = getLogger();

export default function unasyncify(instance: AsyncDependencyExtractor): IDependencyExtractor {
    return {
        match(request, callback) {
            logger.trace(`[service] match request: ${JSON.stringify(request)}`);

            instance.match(request)
                .then((response) => callback(null, response))
                .catch((error) => {
                    logger.error(`[service] match error: ${error}`);
                    callback(error, null);
                });
        },
        extract(request, callback) {
            logger.trace(`[service] extract request: ${JSON.stringify(request)}`);

            instance.extract(request)
                .then((response) => callback(null, response))
                .catch((error) => {
                    logger.error(`[service] extract error: ${error}`);
                    callback(error, null);
                });
        },
    };
}

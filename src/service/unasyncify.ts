import {getLogger} from "log4js";
import { IDependencyExtractor } from "../../api/extractor";
import AsyncDependencyExtractor from "./AsyncDependencyExtractor";
const logger = getLogger();

export default function unasyncify(instance: AsyncDependencyExtractor): IDependencyExtractor {
    return {
        extract(request, callback) {
            logger.info(`[service] extract request: ${JSON.stringify(request)}`);

            instance.extract(request)
                .then((response) => callback(null, response))
                .catch((error) => {
                    logger.info(`[service] extract error: ${error}`);
                    callback(error, null);
                });
        },
    };
}

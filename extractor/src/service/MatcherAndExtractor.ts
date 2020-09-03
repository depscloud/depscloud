import Matcher from "../matcher/Matcher";
import Extractor from "../extractors/Extractor";

export default interface MatcherAndExtractor {
    matcher: Matcher;
    extractor: Extractor;
}

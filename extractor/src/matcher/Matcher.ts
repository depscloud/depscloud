import MatchConfig from "./MatchConfig";

import { Minimatch } from "minimatch";

export default class Matcher {
    private includes: Minimatch[];
    private excludes: Minimatch[];

    constructor(config: MatchConfig) {
        this.includes = config.includes.map((pattern) => new Minimatch(pattern));
        this.excludes = config.excludes.map((pattern) => new Minimatch(pattern));
    }

    public match(path: string): boolean {
        let included = false;
        for (let i = 0; i < this.includes.length && !included; i++) {
            included = this.includes[i].match(path);
        }

        if (!included) {
            return false
        }

        let excluded = false;
        for (let i = 0; i < this.excludes.length && !excluded; i++) {
            excluded = this.excludes[i].match(path);
        }
        return !excluded
    }
}

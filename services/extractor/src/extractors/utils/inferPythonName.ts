import path = require("path");
import url = require("url");

export default function inferPythonName(href: string): string {
    if (!href) {
        return "";
    }

    if (href.indexOf("://") === -1) {
        // assuming scp url

        const [ userAndHost, pathAndExt ] = href.split(":");

        href = "ssh://" + userAndHost + "/" + pathAndExt;
    }

    const parts = url.parse(href);

    return path.basename(parts.path, path.extname(parts.path))
}

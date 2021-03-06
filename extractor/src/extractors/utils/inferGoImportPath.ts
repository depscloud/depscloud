import path = require("path");
import url = require("url");

export default function inferGoImportPath(href: string): string {
    if (!href) {
        return "";
    }

    if (href.indexOf("://") === -1) {
        // assuming scp url

        const [ userAndHost, pathAndExt ] = href.split(":");

        href = "ssh://" + userAndHost + "/" + pathAndExt;
    }

    const parts = url.parse(href);

    const extlength = path.extname(parts.path).length
    return parts.host + parts.path.substr(0, parts.path.length - extlength);
}

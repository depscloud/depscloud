import url = require("url");

export default function inferImportPath(href: string): string {
    if (!href) {
        return "";
    }

    if (href.indexOf("://") === -1) {
        // assuming scp url

        const [ userAndHost, pathAndExt ] = href.split(":");

        href = "ssh://" + userAndHost + "/" + pathAndExt;
    }

    const parts = url.parse(href);

    return parts.host + parts.path.substr(0, parts.path.length - 4);
}

export interface Parsed {
    organization: string;
    module: string;
}

export default function parseImportPath(importPath: string): Parsed {
    const pos = importPath.indexOf("/");
    const organization = importPath.substr(0, pos);
    const module = importPath.substr(pos + 1);
    return { organization, module };
}

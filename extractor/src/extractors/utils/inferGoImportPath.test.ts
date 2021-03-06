import inferGoImportPath from "./inferGoImportPath";

describe("inferImportPath", () => {
    test("git-ssh", () => {
        const importPath = inferGoImportPath("git@github.com:depscloud/depscloud-project.git");

        expect(importPath).toBe("github.com/depscloud/depscloud-project");
    });

    test("git-https", () => {
        const importPath = inferGoImportPath("https://github.com/depscloud/depscloud-project.git");

        expect(importPath).toBe("github.com/depscloud/depscloud-project");
    });
});

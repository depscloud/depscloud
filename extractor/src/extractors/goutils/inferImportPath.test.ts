import inferImportPath from "./inferImportPath";

describe("inferImportPath", () => {
    test("git-ssh", () => {
        const importPath = inferImportPath("git@github.com:depscloud/depscloud-project.git");

        expect(importPath).toBe("github.com/depscloud/depscloud-project");
    });

    test("git-https", () => {
        const importPath = inferImportPath("https://github.com/depscloud/depscloud-project.git");

        expect(importPath).toBe("github.com/depscloud/depscloud-project");
    });
});

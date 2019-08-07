import inferImportPath from "./inferImportPath";

describe("inferImportPath", () => {
    test("git-ssh", () => {
        const importPath = inferImportPath("git@github.com:deps-cloud/deps-cloud-project.git");

        expect(importPath).toBe("github.com/deps-cloud/deps-cloud-project");
    });

    test("git-https", () => {
        const importPath = inferImportPath("https://github.com/deps-cloud/deps-cloud-project.git");

        expect(importPath).toBe("github.com/deps-cloud/deps-cloud-project");
    });
});

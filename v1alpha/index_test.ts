describe("v1alpha", () => {
    test("require", () => {
        const schema = require("./index");

        Object.keys(schema).forEach((key) => {
            console.log(key)
        })
    });
});

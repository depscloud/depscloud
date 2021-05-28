export default class Registry<T> {
    private readonly type: string;
    private readonly registry: { [key: string]: (params: any) => Promise<T> };

    constructor(type: string) {
        this.type = type;
        this.registry = {};
    }

    public registerAll(all: { [key: string]: (params: any) => Promise<T> }): void {
        Object.keys(all)
            .forEach((name) => this.register(name, all[name]));
    }

    public register(name: string, fn: (params: any) => Promise<T>): void {
        this.registry[name] = fn;
    }

    public resolve(name: string, params: any): Promise<T> {
        if (!this.registry[name]) {
            throw new Error(`${this.type} with name: ${name} not registered.`);
        }

        return this.registry[name](params);
    }

    public known(): string[] {
        return Object.keys(this.registry);
    }
}

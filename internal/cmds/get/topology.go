package get

import (
	"fmt"

	"github.com/deps-cloud/api/v1alpha/schema"
	"github.com/deps-cloud/api/v1alpha/tracker"
	"github.com/deps-cloud/cli/internal/writer"

	"github.com/spf13/cobra"
)

func key(module *schema.Module) string {
	return fmt.Sprintf("%s|%s|%s",
		module.Language,
		module.Organization,
		module.Module)
}

type entry struct {
	module *schema.Module
	count  uint64
}

type fetcher func(req *tracker.DependencyRequest) ([]*tracker.Dependency, error)

func topology(root *schema.Module, fetch fetcher) ([][]*schema.Module, error) {
	dependencies := make(map[string][]string)

	idx := make(map[string]*entry)

	current := []*schema.Module{root}
	idx[key(root)] = &entry{
		module: root,
		count:  0,
	}

	for length := len(current); length > 0; length = len(current) {
		next := make([]*schema.Module, 0)

		for i := 0; i < length; i++ {
			module := current[i]

			results, err := fetch(&tracker.DependencyRequest{
				Language:     module.Language,
				Organization: module.Organization,
				Module:       module.Module,
			})
			if err != nil {
				return nil, err
			}

			keys := make([]string, len(results))
			modules := make([]*schema.Module, 0, len(results))

			for i, dependency := range results {
				m := dependency.Module
				k := key(m)

				// always set the key so we decrement later
				keys[i] = k

				// only add the module for processing when we haven't seen it before
				if _, ok := idx[k]; ok {
					idx[k].count++
				} else {
					idx[k] = &entry{
						module: m,
						count:  1,
					}

					modules = append(modules, m)
				}
			}

			dependencies[key(module)] = keys
			next = append(next, modules...)
		}

		current = next
	}

	modules := []string{key(root)}
	result := make([][]*schema.Module, 0)

	for length := len(modules); length > 0; length = len(modules) {
		next := make([]string, 0)
		tier := make([]*schema.Module, 0)

		for i := 0; i < length; i++ {
			k := modules[i]

			results := dependencies[k]
			delete(dependencies, k)

			for _, dependency := range results {
				idx[dependency].count--

				if idx[dependency].count == 0 {
					next = append(next, dependency)
					tier = append(tier, idx[dependency].module)
					delete(idx, dependency)
				}
			}
		}

		modules = next
		if len(tier) > 0 {
			result = append(result, tier)
		}
	}

	return result, nil
}

func TopologyCommand(
	dependencyClient tracker.DependencyServiceClient,
	writer writer.Writer,
) *cobra.Command {
	req := &schema.Module{}

	tiered := false

	cmd := &cobra.Command{
		Use:     "topology <dependents|dependencies>",
		Short:   "Get the module topology of either dependents or dependencies",
		Example: "depscloud-cli get topology dependents -l go -o github.com -m deps-cloud/api",
		RunE: func(cmd *cobra.Command, args []string) error {
			if req.Language == "" || req.Organization == "" || req.Module == "" {
				return fmt.Errorf("language, organization, and module must be provided")
			}

			ctx := cmd.Context()

			var results [][]*schema.Module
			var err error

			switch args[0] {
			case "dependents":
				results, err = topology(req, func(req *tracker.DependencyRequest) ([]*tracker.Dependency, error) {
					resp, err := dependencyClient.ListDependents(ctx, req)
					if err != nil {
						return nil, err
					}
					return resp.Dependents, nil
				})
			case "dependencies":
				results, err = topology(req, func(req *tracker.DependencyRequest) ([]*tracker.Dependency, error) {
					resp, err := dependencyClient.ListDependencies(ctx, req)
					if err != nil {
						return nil, err
					}
					return resp.Dependencies, nil
				})
			default:
				return fmt.Errorf("unrecognized kind: %s", args[0])
			}

			if err != nil {
				return err
			}

			for _, tier := range results {
				if tiered {
					if err := writer.Write(tier); err != nil {
						return err
					}
				} else {
					for _, module := range tier {
						if err := writer.Write(module); err != nil {
							return err
						}
					}
				}
			}

			return nil
		},
	}

	addModuleFlags(cmd, req)

	cmd.Flags().BoolVar(&tiered, "tiered", tiered, "Produce a tiered output instead of a flat stream")

	return cmd
}

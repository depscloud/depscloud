package get

import (
	"context"
	"fmt"

	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/cli/internal/writer"

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
	seen   map[string]bool
}

type fetcher func(req *tracker.DependencyRequest, ctx context.Context) ([]*tracker.Dependency, error)

func topology(root *schema.Module, ctx context.Context, fetch fetcher) ([][]*schema.Module, error) {
	edges := make(map[string][]string)
	nodes := make(map[string]*entry)

	current := []*schema.Module{root}
	nodes[key(root)] = &entry{
		module: root,
		seen:   make(map[string]bool),
	}

	for length := len(current); length > 0; length = len(current) {
		next := make([]*schema.Module, 0)

		for i := 0; i < length; i++ {
			module := current[i]
			moduleKey := key(module)

			results, err := fetch(&tracker.DependencyRequest{
				Language:     module.Language,
				Organization: module.Organization,
				Module:       module.Module,
			}, ctx)
			if err != nil {
				return nil, err
			}

			keys := make([]string, len(results))
			modules := make([]*schema.Module, 0, len(results))

			for i, dependency := range results {
				dependencyKey := key(dependency.Module)

				// always set the key so we decrement later
				keys[i] = dependencyKey

				// only add the module for processing when we haven't seen it before
				if _, ok := nodes[dependencyKey]; !ok {
					nodes[dependencyKey] = &entry{
						module: dependency.Module,
						seen:   make(map[string]bool),
					}

					modules = append(modules, dependency.Module)
				}

				nodes[dependencyKey].seen[moduleKey] = true
			}

			edges[moduleKey] = keys
			next = append(next, modules...)
		}

		current = next
	}

	modules := []string{key(root)}
	result := [][]*schema.Module{{root}}

	for length := len(modules); length > 0; length = len(modules) {
		next := make([]string, 0)
		tier := make([]*schema.Module, 0)

		for i := 0; i < length; i++ {
			k := modules[i]

			results := edges[k]
			delete(edges, k)

			for _, dependencyKey := range results {
				// cycle in the graph, dependencyKey no longer exists
				if _, ok := nodes[dependencyKey]; !ok {
					continue
				}

				delete(nodes[dependencyKey].seen, k)

				if len(nodes[dependencyKey].seen) == 0 {
					next = append(next, dependencyKey)
					tier = append(tier, nodes[dependencyKey].module)
					delete(nodes, dependencyKey)
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

func topologyCommand(
	writer writer.Writer,
	fetch fetcher,
) *cobra.Command {
	req := &schema.Module{}
	tiered := false

	cmd := &cobra.Command{
		Use:   "topology",
		Aliases: []string{"topo"},
		Short: "Get the associated topology",
		RunE: func(cmd *cobra.Command, args []string) error {
			if req.Language == "" || req.Organization == "" || req.Module == "" {
				return fmt.Errorf("language, organization, and module must be provided")
			}

			results, err := topology(req, cmd.Context(), fetch)

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

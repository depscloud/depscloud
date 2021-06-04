package get

import (
	"context"
	"fmt"
	"io"

	"github.com/depscloud/api/v1beta"
	"github.com/depscloud/depscloud/services/deps/internal/writer"

	"github.com/spf13/cobra"
)

func key(module *v1beta.Module) string {
	return fmt.Sprintf("%s|%s", module.Language, module.Name)
}

type convertRequest func(request *v1beta.Module) *v1beta.SearchRequest

func topology(ctx context.Context, searchService v1beta.TraversalServiceClient, request *v1beta.SearchRequest) ([][]*v1beta.Module, error) {
	call, err := searchService.BreadthFirstSearch(ctx)
	if err != nil {
		return nil, err
	}

	if err := call.Send(request); err != nil {
		return nil, err
	}

	nodes := make(map[string]*v1beta.Module)
	counter := make(map[string]map[string]bool)
	edges := make(map[string][]string)

	for resp, err := call.Recv(); true; resp, err = call.Recv() {
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		var source *v1beta.Dependency
		var items []*v1beta.Dependency

		if source = resp.GetRequest().GetDependentsOf(); source != nil {
			items = resp.GetDependents()
		} else if source = resp.GetRequest().GetDependenciesFor(); source != nil {
			items = resp.GetDependencies()
		} else {
			return nil, fmt.Errorf("source module not included")
		}

		moduleKey := key(source.Module)
		if _, ok := nodes[moduleKey]; !ok {
			nodes[moduleKey] = source.Module
			counter[moduleKey] = make(map[string]bool)
		}

		dependencyKeys := make([]string, len(items))

		for i, dependency := range items {
			dependencyKey := key(dependency.Module)

			if _, ok := nodes[dependencyKey]; !ok {
				nodes[dependencyKey] = dependency.Module
				counter[dependencyKey] = make(map[string]bool)
			}

			dependencyKeys[i] = dependencyKey
			counter[dependencyKey][moduleKey] = true
		}

		edges[moduleKey] = dependencyKeys
	}

	var root *v1beta.Dependency
	if root = request.GetDependentsOf(); root != nil {
		// empty
	} else if root = request.GetDependenciesFor(); root != nil {
		// empty
	} else {
		return nil, fmt.Errorf("failed to determine root key for topological sort")
	}

	rootKey := key(root.Module)

	modules := []string{rootKey}
	result := [][]*v1beta.Module{{nodes[rootKey]}}
	delete(nodes, rootKey)
	delete(counter, rootKey)

	for length := len(modules); length > 0; length = len(modules) {
		next := make([]string, 0)
		tier := make([]*v1beta.Module, 0)

		for i := 0; i < length; i++ {
			key := modules[i]

			dependencyKeys := edges[key]
			delete(edges, key)

			for _, dependencyKey := range dependencyKeys {
				// cycle in the graph, dependencyKey no longer exists
				if _, ok := counter[dependencyKey]; !ok {
					continue
				}

				delete(counter[dependencyKey], key)

				if len(counter[dependencyKey]) == 0 {
					next = append(next, dependencyKey)
					tier = append(tier, nodes[dependencyKey])

					delete(nodes, dependencyKey)
					delete(counter, dependencyKey)
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
	traversalService v1beta.TraversalServiceClient,
	requestConverter convertRequest,
) *cobra.Command {
	req := &v1beta.Module{}
	tiered := false

	cmd := &cobra.Command{
		Use:     "tree",
		Aliases: []string{"topology", "topo"},
		Short:   "Get the associated tree using the provided module as the root",
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := topology(cmd.Context(), traversalService, requestConverter(req))

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

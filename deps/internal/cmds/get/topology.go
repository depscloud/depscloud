package get

import (
	"context"
	"fmt"
	"io"

	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/depscloud/deps/internal/writer"

	"github.com/spf13/cobra"
)

func requestToModule(req *tracker.DependencyRequest) *schema.Module {
	return &schema.Module{
		Language:     req.Language,
		Organization: req.Organization,
		Module:       req.Module,
		Name:         req.Name,
	}
}

func key(module *schema.Module) string {
	return fmt.Sprintf("%s|%s|%s|%s",
		module.Language,
		module.Organization,
		module.Module,
		module.Name)
}

func keyForRequest(req *tracker.DependencyRequest) string {
	return fmt.Sprintf("%s|%s|%s|%s",
		req.Language,
		req.Organization,
		req.Module,
		req.Name)
}

type convertRequest func(request *tracker.DependencyRequest) *tracker.SearchRequest

func topology(ctx context.Context, searchService tracker.SearchServiceClient, request *tracker.SearchRequest) ([][]*schema.Module, error) {
	call, err := searchService.BreadthFirstSearch(ctx)
	if err != nil {
		return nil, err
	}

	if err := call.Send(request); err != nil {
		return nil, err
	}

	counter := make(map[string]map[string]bool)
	nodes := make(map[string]*schema.Module)
	edges := make(map[string][]string)

	for resp, err := call.Recv(); true; resp, err = call.Recv() {
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		var source *tracker.DependencyRequest
		var items []*tracker.Dependency

		if source = resp.GetRequest().GetDependentsOf(); source != nil {
			items = resp.GetDependents()
		} else if source = resp.GetRequest().GetDependenciesOf(); source != nil {
			items = resp.GetDependencies()
		}

		if source == nil {
			return nil, fmt.Errorf("source module not included")
		}

		module := requestToModule(source)
		moduleKey := key(module)
		if _, ok := nodes[moduleKey]; !ok {
			nodes[moduleKey] = module
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

	var rootKey string

	if req := request.GetDependentsOf(); req != nil {
		rootKey = keyForRequest(req)
	} else if req := request.GetDependenciesOf(); req != nil {
		rootKey = keyForRequest(req)
	} else {
		return nil, fmt.Errorf("failed to determine root key for topological sort")
	}

	modules := []string{rootKey}
	result := [][]*schema.Module{{nodes[rootKey]}}
	delete(nodes, rootKey)
	delete(counter, rootKey)

	for length := len(modules); length > 0; length = len(modules) {
		next := make([]string, 0)
		tier := make([]*schema.Module, 0)

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
	searchService tracker.SearchServiceClient,
	requestConverter convertRequest,
) *cobra.Command {
	req := &tracker.DependencyRequest{}
	tiered := false

	cmd := &cobra.Command{
		Use:     "topology",
		Aliases: []string{"topo"},
		Short:   "Get the associated topology",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateDependencyRequest(req); err != nil {
				return err
			}

			results, err := topology(cmd.Context(), searchService, requestConverter(req))

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

	addDependencyRequestFlags(cmd, req)

	cmd.Flags().BoolVar(&tiered, "tiered", tiered, "Produce a tiered output instead of a flat stream")

	return cmd
}

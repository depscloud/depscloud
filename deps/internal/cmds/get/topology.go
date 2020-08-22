package get

import (
	"context"
	"fmt"
	"io"

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

func keyForRequest(req *tracker.DependencyRequest) string {
	return fmt.Sprintf("%s|%s|%s",
		req.Language,
		req.Organization,
		req.Module)
}

type entry struct {
	req    *tracker.DependencyRequest
	module *schema.Module
	seen   map[string]bool
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

	edges := make(map[string][]string)
	nodes := make(map[string]*entry)

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

		sourceKey := keyForRequest(source)
		targets := make([]string, len(items))

		for i, dependency := range items {
			targets[i] = key(dependency.Module)

			if _, ok := nodes[targets[i]]; !ok {
				nodes[targets[i]] = &entry{
					module: dependency.GetModule(),
					seen: make(map[string]bool),
				}
			}

			nodes[targets[i]].seen[sourceKey] = true
		}

		edges[sourceKey] = targets
	}

	var root *schema.Module

	if req := request.GetDependentsOf(); req != nil {
		root = &schema.Module{
			Language: req.GetLanguage(),
			Organization: req.GetOrganization(),
			Module: req.GetModule(),
		}
	} else if req := request.GetDependenciesOf(); req != nil {
		root = &schema.Module{
			Language: req.GetLanguage(),
			Organization: req.GetOrganization(),
			Module: req.GetModule(),
		}
	}

	if root == nil {
		return nil, fmt.Errorf("failed to determine root key for topological sort")
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
	searchService tracker.SearchServiceClient,
	requestConverter convertRequest,
) *cobra.Command {
	req := &tracker.DependencyRequest{}
	tiered := false

	cmd := &cobra.Command{
		Use:   "topology",
		Aliases: []string{"topo"},
		Short: "Get the associated topology",
		RunE: func(cmd *cobra.Command, args []string) error {
			if req.Language == "" || req.Organization == "" || req.Module == "" {
				return fmt.Errorf("language, organization, and module must be provided")
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

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

func getAllPaths(ctx context.Context, searchService tracker.SearchServiceClient, request *tracker.SearchRequest, destinationNodeKey string) ([][]*schema.Module, error) {
	// Unclear as to what's the significant advantage we are achieving with using DFS over BFS
	// As it's written right now, we do wait until the search is fully exhausted (the for loop below) anyway,
	// before we begin to build the topology
	call, err := searchService.DepthFirstSearch(ctx)
	if err != nil {
		return nil, err
	}

	if err := call.Send(request); err != nil {
		return nil, err
	}
	nodes := make(map[string]*schema.Module) // key: Module, value: Metadata of the module
	edges := make(map[string][]string)       // key: Module, value: adjacency list of the module

	// loop until all the responses are received from the server
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
		}

		dependencyKeys := make([]string, len(items))

		for i, dependency := range items {
			dependencyKey := key(dependency.Module)

			if _, ok := nodes[dependencyKey]; !ok {
				nodes[dependencyKey] = dependency.Module
			}

			dependencyKeys[i] = dependencyKey
		}

		edges[moduleKey] = dependencyKeys
	}

	var rootKey string
	if req := request.GetDependentsOf(); req != nil {
		rootKey = keyForRequest(req)
	} else if req := request.GetDependenciesOf(); req != nil {
		rootKey = keyForRequest(req)
	} else {
		return nil, fmt.Errorf("failed to determine root key for topological paths")
	}

	stack := [][]string{{rootKey}}
	result := [][]*schema.Module{{}}

	for length := len(stack); length > 0; length = len(stack) {
		next := make([][]string, 0)

		// Pop
		path := stack[length-1]
		stack = stack[0:(length - 1)]

		if len(path) == 0 {
			continue // a path in the stack should ideally never be empty
		}

		// Continue exploration of this path
		node := path[len(path)-1]
		if node == destinationNodeKey {
			// We've found a matching path
			var modules []*schema.Module
			for _, moduleKey := range path {
				// Translate keys to metadata of the module
				modules = append(modules, nodes[moduleKey])
			}
			result = append(result, modules)

			continue
		}

		// Append new edges to the path to continue the search for destination node
		nodeEdges := edges[node]
		for _, newEdge := range nodeEdges {
			edgeExists := false
			// Check and avoid cycles
			// TODO: Can we do better than this? Currently, it iterates over the path for every new edge
			for _, edgeInPath := range path {
				if edgeInPath == newEdge {
					edgeExists = true
					break
				}
			}

			if !edgeExists {
				// Copy "path" and append, so we push only "distinct" paths to the stack
				pathCopy := make([]string, len(path), cap(path)+1)
				copy(pathCopy, path)
				next = append(next, append(pathCopy, newEdge))
			}
		}
		stack = append(stack, next...)
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
	var destinationModuleName string

	cmd := &cobra.Command{
		Use:     "topology",
		Aliases: []string{"topo"},
		Short:   "Get the associated topology",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateDependencyRequest(req); err != nil {
				return err
			}

			if !(destinationModuleName == "") {
				orgAndModule := parseName(req.GetLanguage(), destinationModuleName)
				destinationModuleKey := key(&schema.Module{
					Language:     req.GetLanguage(),
					Organization: orgAndModule[0],
					Module:       orgAndModule[1],
					Name:         destinationModuleName,
				})

				results, err := getAllPaths(cmd.Context(), searchService, requestConverter(req), destinationModuleKey)
				if err != nil {
					return err
				}
				for _, paths := range results {
					if err := writer.Write(paths); err != nil {
						return err
					}
				}

				return nil
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
	cmd.Flags().StringVarP(&destinationModuleName, "destinationModuleName", "", "dest", "Get all the ways the module provided with --name flag is connected to this destination module")

	return cmd
}

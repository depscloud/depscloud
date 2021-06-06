package v1beta

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"github.com/depscloud/api/v1beta"
	"github.com/depscloud/api/v1beta/graphstore"
	"github.com/depscloud/depscloud/internal/logger"

	"go.uber.org/zap"

	"google.golang.org/grpc"
)

// TODO: move this to internal to be shared between indexer and tracker
type DefaultKind = string

const (
	ProviderDefaultKind   DefaultKind = "provider"
	RepositoryDefaultKind DefaultKind = "repository"
	ArtifactDefaultKind   DefaultKind = "artifact"
)

func RegisterManifestStorageServiceServer(server *grpc.Server, graphStore graphstore.GraphStoreClient, index IndexService) {
	v1beta.RegisterManifestStorageServiceServer(server, &manifestStorageService{
		graphStore: graphStore,
		index:      index,
	})
}

type manifestStorageService struct {
	v1beta.UnsafeManifestStorageServiceServer

	graphStore graphstore.GraphStoreClient
	index      IndexService
}

func (m *manifestStorageService) GetProposed(request *v1beta.StoreRequest) ([]*graphstore.Node, []*graphstore.Edge, []*Index) {
	nodes := make([]*graphstore.Node, 0)
	edges := make([]*graphstore.Edge, 0)
	index := make([]*Index, 0)

	source, _ := newNode(&v1beta.Source{
		Kind: request.GetKind(),
		Url:  request.GetUrl(),
	})

	ref := request.GetUrl()

	for _, manifestFile := range request.GetManifestFiles() {
		language := manifestFile.GetLanguage()
		system := manifestFile.GetSystem()
		version := manifestFile.GetVersion()

		module, _ := newNode(&v1beta.Module{
			Language: language,
			Name:     manifestFile.GetName(),
		})

		sm := &v1beta.SourceModule{
			Version: version,
			System:  system,
		}

		sourceModule, _ := newEdge(sm)
		sourceModule.FromKey = source.Key
		sourceModule.ToKey = module.Key

		for _, manifestDependency := range manifestFile.GetDependencies() {
			dependency, _ := newNode(&v1beta.Module{
				Language: language,
				Name:     manifestDependency.GetName(),
			})

			moduleDependency, _ := newEdge(&v1beta.ModuleDependency{
				Ref:               ref,
				VersionConstraint: manifestDependency.GetVersionConstraint(),
				Scopes:            manifestDependency.GetScopes(),
			})

			moduleDependency.FromKey = module.Key
			moduleDependency.ToKey = dependency.Key
			moduleDependency.Key = source.Key

			nodes = append(nodes, dependency)
			edges = append(edges, moduleDependency)
			index = append(index,
				&Index{
					Kind:  dependency.GetBody().GetTypeUrl(),
					Field: "language",
					Value: language,
					Key:   base64.StdEncoding.EncodeToString(dependency.Key),
				}, &Index{
					Kind:  dependency.GetBody().GetTypeUrl(),
					Field: "name",
					Value: manifestDependency.GetName(),
					Key:   base64.StdEncoding.EncodeToString(dependency.Key),
				},
			)
		}

		nodes = append(nodes, module)
		edges = append(edges, sourceModule)
		index = append(index,
			&Index{
				Kind:  module.GetBody().GetTypeUrl(),
				Field: "language",
				Value: language,
				Key:   base64.StdEncoding.EncodeToString(module.Key),
			}, &Index{
				Kind:  module.GetBody().GetTypeUrl(),
				Field: "name",
				Value: manifestFile.GetName(),
				Key:   base64.StdEncoding.EncodeToString(module.Key),
			},
		)

		// look for a reported source url

		if sourceURL := manifestFile.GetSourceUrl(); sourceURL != "" {
			reportedSource, _ := newNode(&v1beta.Source{
				Kind: RepositoryDefaultKind,
				Url:  sourceURL,
			})

			reportedSourceModule, _ := newEdge(sm)
			reportedSourceModule.FromKey = reportedSource.Key
			reportedSourceModule.ToKey = module.Key

			nodes = append(nodes, reportedSource)
			edges = append(edges, reportedSourceModule)
			index = append(index,
				&Index{
					Kind:  reportedSource.GetBody().GetTypeUrl(),
					Field: "url",
					Value: sourceURL,
					Key:   base64.StdEncoding.EncodeToString(reportedSource.Key),
				},
			)
		}
	}

	nodes = append(nodes, source)
	index = append(index,
		&Index{
			Kind:  source.GetBody().GetTypeUrl(),
			Field: "url",
			Value: ref,
			Key:   base64.StdEncoding.EncodeToString(source.Key),
		},
	)

	return nodes, edges, index
}

func (m *manifestStorageService) GetStored(ctx context.Context, source *graphstore.Node) ([]*graphstore.Node, []*graphstore.Edge, error) {
	call, err := m.graphStore.Traverse(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer call.SendMsg(&graphstore.TraverseRequest{
		Cancel: true,
	})

	tier := map[string]*graphstore.Node{
		string(source.Key): source,
	}

	nodes := []*graphstore.Node{source}
	edges := make([]*graphstore.Edge, 0)

	// (Source) => SourceModule => (Module) => ModuleDependency => (Module)
	//    0                           1
	for depth := 0; depth < 2; depth++ {
		for _, node := range tier {
			err = call.SendMsg(&graphstore.TraverseRequest{
				Request: &graphstore.NeighborsRequest{
					From: node,
				},
			})
			if err != nil {
				return nil, nil, err
			}
		}

		next := make(map[string]*graphstore.Node)

		// all requests sent, await all responses
		for i := 0; i < len(tier); i++ {
			resp, err := call.Recv()
			if err != nil {
				return nil, nil, err
			}

			for _, neighbor := range resp.Response.Neighbors {
				next[string(neighbor.Node.Key)] = neighbor.Node

				nodes = append(nodes, neighbor.Node)
				for _, edge := range neighbor.Edges {
					if len(edge.Key) == 0 || bytes.Equal(edge.Key, source.Key) {
						edges = append(edges, edge)
					}
				}
			}
		}

		tier = next
	}

	return nodes, edges, nil
}

func (m *manifestStorageService) Store(ctx context.Context, request *v1beta.StoreRequest) (*v1beta.StoreResponse, error) {
	log := logger.Extract(ctx)

	proposedNodes, proposedEdges, indexFields := m.GetProposed(request)

	// last node is the provided source
	source := proposedNodes[len(proposedNodes)-1]

	// we don't care about stored nodes in this process because we don't ever delete nodes
	_, storedEdges, err := m.GetStored(ctx, source)
	if err != nil {
		log.Error(ErrQueryFailure.Error(), zap.Error(err))
		return nil, ErrQueryFailure
	}

	edgeIndex := make(map[string]*graphstore.Edge, len(proposedEdges))
	for _, edge := range proposedEdges {
		key := fmt.Sprintf("%s/%s/%s",
			string(edge.FromKey), string(edge.ToKey), string(edge.Key))

		edgeIndex[key] = edge
	}

	edgesToDelete := make([]*graphstore.Edge, 0)
	for _, edge := range storedEdges {
		key := fmt.Sprintf("%s/%s/%s",
			string(edge.FromKey), string(edge.ToKey), string(edge.Key))

		if _, ok := edgeIndex[key]; !ok {
			// not in the proposed set, needs to be removed
			edgesToDelete = append(edgesToDelete, edge)
		}
	}

	_, err = m.graphStore.Put(ctx, &graphstore.PutRequest{
		Nodes: proposedNodes,
		Edges: proposedEdges,
	})
	if err != nil {
		log.Error("failed to update graph with new data", zap.Error(err))
		return nil, ErrUpdateFailure
	}

	_, err = m.graphStore.Delete(ctx, &graphstore.DeleteRequest{
		Edges: edgesToDelete,
	})
	if err != nil {
		log.Error("failed to remove outdated edges", zap.Error(err))
		return nil, ErrPruneFailure
	}

	err = m.index.Index(ctx, indexFields)

	return &v1beta.StoreResponse{}, err
}

var _ v1beta.ManifestStorageServiceServer = &manifestStorageService{}

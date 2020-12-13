package v1alpha

import (
	"bytes"
	"context"

	"github.com/depscloud/api"
	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/store"
	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/depscloud/internal/logger"

	"go.uber.org/zap"

	"google.golang.org/grpc"
)

// RegisterSourceService registers the sourceService implementation with the server
func RegisterSourceService(server *grpc.Server, gs store.GraphStoreClient) {
	tracker.RegisterSourceServiceServer(server, &sourceService{gs: gs})
}

type sourceService struct {
	gs store.GraphStoreClient
}

var _ tracker.SourceServiceServer = &sourceService{}

func (s *sourceService) List(ctx context.Context, req *tracker.ListRequest) (*tracker.ListSourceResponse, error) {
	log := logger.Extract(ctx)

	resp, err := s.gs.List(ctx, &store.ListRequest{
		Page:  req.GetPage(),
		Count: req.GetCount(),
		Type:  SourceType,
	})

	if err != nil {
		log.Error("failed to list sources", zap.Error(err))
		return nil, err
	}

	sources := make([]*schema.Source, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		source, _ := Decode(item)
		sources = append(sources, source.(*schema.Source))
	}

	return &tracker.ListSourceResponse{
		Page:    req.GetPage(),
		Count:   req.GetCount(),
		Sources: sources,
	}, nil
}

func (s *sourceService) Track(ctx context.Context, req *tracker.SourceRequest) (*tracker.TrackResponse, error) {
	source := req.GetSource()
	log := logger.Extract(ctx).With(
		zap.String("source_url", source.GetUrl()),
		zap.String("source_kind", source.GetKind()),
		zap.String("source_ref", source.GetRef()))

	currentSet, err := s.getCurrent(ctx, source)
	if err != nil {
		log.Error("failed to get current tree", zap.Error(err))
		return nil, api.ErrModuleNotFound
	}

	proposedSet, err := s.getProposed(ctx, req)
	if err != nil {
		log.Error("failed to get proposed tree", zap.Error(err))
		return nil, api.ErrModuleNotFound
	}

	toDelete := make([]*store.GraphItem, 0)
	for key, item := range currentSet {
		if _, ok := proposedSet[key]; !ok {
			// don't delete modules from the graph when an edge to it is removed.
			// we'll put a cleanup in later
			if item.GetGraphItemType() != ModuleType {
				toDelete = append(toDelete, item)
			}
		}
	}

	toPut := make([]*store.GraphItem, 0, len(proposedSet))
	for _, item := range proposedSet {
		toPut = append(toPut, item)
	}

	log.Info("proposed changes",
		zap.Int("current_set", len(currentSet)),
		zap.Int("proposed_set", len(proposedSet)),
		zap.Int("to_put", len(toPut)),
		zap.Int("to_delete", len(toDelete)))

	if _, err := s.gs.Delete(ctx, &store.DeleteRequest{Items: toDelete}); err != nil {
		log.Error("failed to delete edge data", zap.Error(err))
		return nil, api.ErrPartialDeletion
	}

	if _, err := s.gs.Put(ctx, &store.PutRequest{Items: toPut}); err != nil {
		log.Error("failed to put node and edge data", zap.Error(err))
		return nil, api.ErrPartialInsertion
	}

	return &tracker.TrackResponse{Tracking: true}, nil
}

func (s *sourceService) getCurrent(ctx context.Context, source *schema.Source) (map[string]*store.GraphItem, error) {
	idx := make(map[string]*store.GraphItem)

	item, err := Encode(source)
	if err != nil {
		return nil, err
	}

	idx[readableKey(item)] = item

	manages, err := s.gs.FindUpstream(ctx, &store.FindRequest{
		Keys:      [][]byte{keyForSource(source)},
		EdgeTypes: []string{ManagesType},
		NodeTypes: []string{ModuleType},
	})

	if err != nil {
		return nil, err
	}

	sourceKey := item.GetK1()
	for _, managed := range manages.GetPairs() {
		idx[readableKey(managed.GetNode())] = managed.GetNode()
		idx[readableKey(managed.GetEdge())] = managed.GetEdge()

		depends, err := s.gs.FindUpstream(ctx, &store.FindRequest{
			Keys:      [][]byte{managed.GetNode().GetK1()},
			EdgeTypes: []string{DependsType},
			NodeTypes: []string{ModuleType},
		})

		if err != nil {
			return nil, err
		}

		for _, depended := range depends.GetPairs() {
			// Return only the depends edges that are produced by modules of this source URL
			dependsEdgeK3 := depended.GetEdge().GetK3()
			if len(dependsEdgeK3) == 0 || bytes.Equal(dependsEdgeK3, sourceKey) {
				idx[readableKey(depended.GetNode())] = depended.GetNode()
				idx[readableKey(depended.GetEdge())] = depended.GetEdge()
			}
		}
	}

	return idx, nil
}

func (s *sourceService) getProposed(ctx context.Context, request *tracker.SourceRequest) (map[string]*store.GraphItem, error) {
	idx := make(map[string]*store.GraphItem)

	source, err := Encode(request.GetSource())
	if err != nil {
		return nil, err
	}

	idx[readableKey(source)] = source

	for _, managementFile := range request.GetManagementFiles() {
		managedModule, err := Encode(&schema.Module{
			Language:     managementFile.GetLanguage(),
			Organization: managementFile.GetOrganization(),
			Module:       managementFile.GetModule(),
			Name:         managementFile.GetName(),
		})

		if err != nil {
			return nil, err
		}

		manages, err := Encode(&schema.Manages{
			Language: managementFile.GetLanguage(),
			System:   managementFile.GetSystem(),
			Version:  managementFile.GetVersion(),
		})
		if err != nil {
			return nil, err
		}

		manages.K1 = source.GetK1()
		manages.K2 = managedModule.GetK1()

		idx[readableKey(managedModule)] = managedModule
		idx[readableKey(manages)] = manages

		if sourceURL := managementFile.GetSourceUrl(); sourceURL != "" {
			discoveredSource, err := Encode(&schema.Source{
				Url:  sourceURL,
				Kind: "repository",
			})
			if err != nil {
				return nil, err
			}

			discoveredManages, err := Encode(&schema.Manages{
				Language: managementFile.GetLanguage(),
				System:   managementFile.GetSystem(),
				Version:  managementFile.GetVersion(),
			})
			if err != nil {
				return nil, err
			}

			discoveredManages.K1 = discoveredSource.GetK1()
			discoveredManages.K2 = managedModule.GetK1()

			idx[readableKey(discoveredSource)] = discoveredSource
			idx[readableKey(discoveredManages)] = discoveredManages
		}

		for _, dependency := range managementFile.GetDependencies() {
			dependedModule, err := Encode(&schema.Module{
				Language:     managementFile.GetLanguage(),
				Organization: dependency.GetOrganization(),
				Module:       dependency.GetModule(),
				Name:         dependency.GetName(),
			})
			if err != nil {
				return nil, err
			}

			depends, err := Encode(&schema.Depends{
				Language:          managementFile.GetLanguage(),
				VersionConstraint: dependency.GetVersionConstraint(),
				Scopes:            dependency.GetScopes(),
				Ref:               request.GetSource().Ref,
			})
			if err != nil {
				return nil, err
			}

			depends.K1 = managedModule.GetK1()
			depends.K2 = dependedModule.GetK1()
			depends.K3 = source.GetK1()

			idx[readableKey(dependedModule)] = dependedModule
			idx[readableKey(depends)] = depends
		}
	}

	return idx, nil
}

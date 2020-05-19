package services

import (
	"bytes"
	"context"

	"github.com/deps-cloud/api"
	"github.com/deps-cloud/api/v1alpha/schema"
	"github.com/deps-cloud/api/v1alpha/store"
	"github.com/deps-cloud/api/v1alpha/tracker"
	"github.com/deps-cloud/tracker/pkg/types"

	"github.com/sirupsen/logrus"

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
	resp, err := s.gs.List(ctx, &store.ListRequest{
		Page:  req.GetPage(),
		Count: req.GetCount(),
		Type:  types.SourceType,
	})

	if err != nil {
		logrus.Errorf("[service.source] %s", err.Error())
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
	currentSet, err := s.getCurrent(ctx, req.GetSource())
	if err != nil {
		logrus.Errorf("[service.source] %s", err.Error())
		return nil, api.ErrModuleNotFound
	}

	proposedSet, err := s.getProposed(ctx, req)
	if err != nil {
		logrus.Errorf("[service.source] %s", err.Error())
		return nil, api.ErrModuleNotFound
	}

	toDelete := make([]*store.GraphItem, 0)
	for key, item := range currentSet {
		if _, ok := proposedSet[key]; !ok {
			toDelete = append(toDelete, item)
		}
	}

	toPut := make([]*store.GraphItem, 0, len(proposedSet))
	for _, item := range proposedSet {
		toPut = append(toPut, item)
	}

	logrus.Infof("[service.source] currentSet=%d proposedSet=%d toDelete=%d toPut=%d",
		len(currentSet), len(proposedSet), len(toDelete), len(toPut))

	if _, err := s.gs.Delete(ctx, &store.DeleteRequest{Items: toDelete}); err != nil {
		logrus.Errorf("[service.source] %s", err.Error())
		return nil, api.ErrPartialDeletion
	}

	if _, err := s.gs.Put(ctx, &store.PutRequest{Items: toPut}); err != nil {
		logrus.Errorf("[service.source] %s", err.Error())
		return nil, api.ErrPartialInsertion
	}

	return &tracker.TrackResponse{Tracking: true}, nil
}

func (s *sourceService) getCurrent(ctx context.Context, source *schema.Source) (map[string]*store.GraphItem, error) {
	idx := make(map[string]*store.GraphItem)

	item, err := Encode(source)
	if err != nil {
		logrus.Errorf("[service.source] %s", err.Error())
		return nil, err
	}

	idx[readableKey(item)] = item

	manages, err := s.gs.FindUpstream(ctx, &store.FindRequest{
		Key:       keyForSource(source),
		EdgeTypes: []string{types.ManagesType},
	})

	if err != nil {
		logrus.Errorf("[service.source] %s", err.Error())
		return nil, err
	}

	sourceKey := item.GetK1()
	for _, managed := range manages.GetPairs() {
		idx[readableKey(managed.GetNode())] = managed.GetNode()
		idx[readableKey(managed.GetEdge())] = managed.GetEdge()

		depends, err := s.gs.FindUpstream(ctx, &store.FindRequest{
			Key:       managed.GetNode().GetK1(),
			EdgeTypes: []string{types.DependsType},
		})

		if err != nil {
			logrus.Errorf("[service.source] %s", err.Error())
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
		logrus.Errorf("[service.source] %s", err.Error())
		return nil, err
	}

	idx[readableKey(source)] = source

	for _, managementFile := range request.GetManagementFiles() {
		managedModule, err := Encode(&schema.Module{
			Language:     managementFile.GetLanguage(),
			Organization: managementFile.GetOrganization(),
			Module:       managementFile.GetModule(),
		})

		if err != nil {
			logrus.Errorf("[service.source] %s", err.Error())
			return nil, err
		}

		manages, err := Encode(&schema.Manages{
			Language: managementFile.GetLanguage(),
			System:   managementFile.GetSystem(),
			Version:  managementFile.GetVersion(),
		})
		if err != nil {
			logrus.Errorf("[service.source] %s", err.Error())
			return nil, err
		}

		manages.K1 = source.GetK1()
		manages.K2 = managedModule.GetK1()

		idx[readableKey(managedModule)] = managedModule
		idx[readableKey(manages)] = manages

		for _, dependency := range managementFile.GetDependencies() {
			dependedModule, err := Encode(&schema.Module{
				Language:     managementFile.GetLanguage(),
				Organization: dependency.GetOrganization(),
				Module:       dependency.GetModule(),
			})
			if err != nil {
				logrus.Errorf("[service.source] %s", err.Error())
				return nil, err
			}

			depends, err := Encode(&schema.Depends{
				Language:          managementFile.GetLanguage(),
				VersionConstraint: dependency.GetVersionConstraint(),
				Scopes:            dependency.GetScopes(),
			})
			if err != nil {
				logrus.Errorf("[service.source] %s", err.Error())
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

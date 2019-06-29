package services

import (
	"context"

	"github.com/deps-cloud/dts/api"
	"github.com/deps-cloud/dts/api/v1alpha"
	"github.com/deps-cloud/dts/api/v1alpha/schema"
	"github.com/deps-cloud/dts/api/v1alpha/store"
	"github.com/deps-cloud/dts/pkg/types"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

// RegisterDependencyService registers the dependencyService implementation with the server
func RegisterDependencyService(server *grpc.Server, gs store.GraphStoreClient) {
	v1alpha.RegisterDependencyServiceServer(server, &dependencyService{gs: gs})
}

type dependencyService struct {
	gs store.GraphStoreClient
}

var _ v1alpha.DependencyServiceServer = &dependencyService{}

func keyForDependencyRequest(req *v1alpha.DependencyRequest) []byte {
	return keyForModule(&schema.Module{
		Language:     req.GetLanguage(),
		Organization: req.GetOrganization(),
		Module:       req.GetModule(),
	})
}

func (d *dependencyService) GetDependents(req *v1alpha.DependencyRequest, resp v1alpha.DependencyService_GetDependentsServer) error {
	key := keyForDependencyRequest(req)

	response, err := d.gs.FindDownstream(context.Background(), &store.FindRequest{
		Key:       key,
		EdgeTypes: []string{types.DependsType},
	})

	if err != nil {
		return api.ErrModuleNotFound
	}

	for _, pair := range response.GetPairs() {
		a, _ := Decode(pair.Node)
		b, _ := Decode(pair.Edge)

		dependency := &v1alpha.Dependency{
			Module:  a.(*schema.Module),
			Depends: b.(*schema.Depends),
		}

		if err := resp.Send(dependency); err != nil {
			logrus.Errorf("[service.dependency] failed to send response: %v", err)
		}
	}

	return nil
}

func (d *dependencyService) GetDependencies(req *v1alpha.DependencyRequest, resp v1alpha.DependencyService_GetDependenciesServer) error {
	key := keyForDependencyRequest(req)

	response, err := d.gs.FindUpstream(context.Background(), &store.FindRequest{
		Key:       key,
		EdgeTypes: []string{types.DependsType},
	})

	if err != nil {
		return api.ErrModuleNotFound
	}

	for _, pair := range response.GetPairs() {
		a, _ := Decode(pair.Node)
		b, _ := Decode(pair.Edge)

		dependency := &v1alpha.Dependency{
			Module:  a.(*schema.Module),
			Depends: b.(*schema.Depends),
		}

		if err := resp.Send(dependency); err != nil {
			logrus.Errorf("[service.dependency] failed to send response: %v", err)
		}
	}

	return nil
}

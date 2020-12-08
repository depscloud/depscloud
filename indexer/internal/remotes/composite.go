package remotes

import (
	"github.com/depscloud/depscloud/internal/logger"

	"go.uber.org/zap"
)

// NewCompositeRemote wraps the supplied remotes in a composite wrapper
// which logs errors and continues processing remote endpoints.
func NewCompositeRemote(remotes ...Remote) Remote {
	return &compositeRemote{
		remotes: remotes,
	}
}

var _ Remote = &compositeRemote{}

type compositeRemote struct {
	remotes []Remote
}

func (r *compositeRemote) FetchRepositories(req *FetchRepositoriesRequest) (*FetchRepositoriesResponse, error) {
	log := logger.Extract(req.Context)

	repositories := make([]*Repository, 0)
	for _, remote := range r.remotes {
		repos, err := remote.FetchRepositories(req)

		if err != nil {
			log.Error("failed to list repositories from remote", zap.Error(err))
		} else {
			repositories = append(repositories, repos.Repositories...)
		}
	}
	return &FetchRepositoriesResponse{
		Repositories: repositories,
	}, nil
}

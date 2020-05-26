package remotes

import "github.com/sirupsen/logrus"

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

func (r *compositeRemote) FetchRepositories(request *FetchRepositoriesRequest) (*FetchRepositoriesResponse, error) {
	repositories := make([]*Repository, 0)
	for _, remote := range r.remotes {
		repos, err := remote.FetchRepositories(request)

		if err != nil {
			logrus.Errorf("[remotes.composite] failed to list repositories from remote: %v", err)
		} else {
			repositories = append(repositories, repos.Repositories...)
		}
	}
	return &FetchRepositoriesResponse{
		Repositories: repositories,
	}, nil
}

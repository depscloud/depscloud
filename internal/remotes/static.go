package remotes

import "github.com/deps-cloud/indexer/internal/config"

var _ Remote = &staticRemote{}

// NewStaticRemote produces a new remote from static configuration
func NewStaticRemote(cfg *config.Static) Remote {
	return &staticRemote{
		config: cfg,
	}
}

type staticRemote struct {
	config *config.Static
}

func (s *staticRemote) FetchRepositories(request *FetchRepositoriesRequest) (*FetchRepositoriesResponse, error) {
	cloneConfig := s.config.GetClone()

	repositories := make([]*Repository, len(s.config.RepositoryUrls))
	for i, repositoryURL := range s.config.RepositoryUrls {
		repositories[i] = &Repository{
			RepositoryURL: repositoryURL,
			Clone:         cloneConfig,
		}
	}
	return &FetchRepositoriesResponse{
		Repositories: repositories,
	}, nil
}

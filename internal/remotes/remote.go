package remotes

import "github.com/deps-cloud/indexer/internal/config"

// Repository represents the combination of a URL and it's corresponding clone credentials.
type Repository struct {
	RepositoryURL string
	Clone         *config.Clone
}

// FetchRepositoriesRequest is a request wrapper that encapsulates request data.
// Currently unused but may be leveraged for filters later on.
type FetchRepositoriesRequest struct {
}

// FetchRepositoriesResponse is a response wrapper that encapsulates response data.
type FetchRepositoriesResponse struct {
	Repositories []*Repository
}

// Remote defines an abstraction for interacting with upstream source control providers.
type Remote interface {
	FetchRepositories(request *FetchRepositoriesRequest) (*FetchRepositoriesResponse, error)
}

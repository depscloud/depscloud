package remotes

import (
	"fmt"

	"github.com/davidji99/bitbucket-go/bitbucket"

	"github.com/depscloud/indexer/internal/config"

	"github.com/sirupsen/logrus"
)

// NewBitbucketRemote constructs a new remote implementation that speaks with Bitbucket
// for repository related information.
func NewBitbucketRemote(cfg *config.Bitbucket) (Remote, error) {
	var client *bitbucket.Client

	if basic := cfg.GetBasic(); basic != nil {
		username := basic.GetUsername()
		password := basic.GetPassword()

		client = bitbucket.NewBasicAuth(username, password)
	} else {
		return nil, fmt.Errorf("auth format not supported")
	}

	return &bitbucketRemote{
		client: client,
		config: cfg,
	}, nil
}

var _ Remote = &bitbucketRemote{}

func convertRepositoriesResponse(response interface{}, cloneConfig *config.Clone) []*Repository {
	rdata := response.(map[string]interface{})
	values := rdata["values"].([]interface{})

	strategy := "ssh"
	if cloneConfig.GetStrategy() == config.CloneStrategy_HTTP {
		strategy = "https"
	}

	repos := make([]*Repository, 0, len(values))
	for _, value := range values {
		val := value.(map[string]interface{})

		links := val["links"].(map[string]interface{})

		cloneURLs := links["clone"].([]interface{})

		for _, cloneURL := range cloneURLs {
			curl := cloneURL.(map[string]interface{})

			if curl["name"].(string) == strategy {
				repos = append(repos, &Repository{
					RepositoryURL: curl["href"].(string),
					Clone:         cloneConfig,
				})
			}
		}
	}

	return repos
}

type bitbucketRemote struct {
	client *bitbucket.Client
	config *config.Bitbucket
}

func (r *bitbucketRemote) FetchRepositories(request *FetchRepositoriesRequest) (*FetchRepositoriesResponse, error) {
	pageLen := uint64(10)
	allRepos := make([]*Repository, 0)
	cloneConfig := r.config.GetClone()

	// if clone config is nil, fall back
	if cloneConfig == nil {
		cloneConfig = &config.Clone{
			Strategy: r.config.GetStrategy(),
		}
	}

	for _, user := range r.config.Users {
		logrus.Infof("[remotes.bitbucket] fetching projects for user: %s", user)

		for page := uint64(1); true; page++ {
			repos, _, err := r.client.Repositories.List(user, &bitbucket.ListOpts{
				Page:    int64(page),
				Pagelen: int64(pageLen),
			})

			if err != nil {
				logrus.Errorf("[remotes.bitbucket] encountered err while fetching projects for user %s, %v", user, err)
				continue
			}

			rr := convertRepositoriesResponse(repos, cloneConfig)
			allRepos = append(allRepos, rr...)

			if uint64(len(rr)) < pageLen {
				break
			}
		}
	}

	for _, team := range r.config.Teams {
		logrus.Infof("[remotes.bitbucket] fetching projects for team: %s", team)

		for page := uint64(1); true; page++ {
			repos, _, err := r.client.Repositories.List(team, &bitbucket.ListOpts{
				Page:    int64(page),
				Pagelen: int64(pageLen),
			})

			if err != nil {
				logrus.Errorf("[remotes.bitbucket] encountered err while fetching projects for team %s, %v", team, err)
				continue
			}

			rr := convertRepositoriesResponse(repos, cloneConfig)
			allRepos = append(allRepos, rr...)

			if uint64(len(rr)) < pageLen {
				break
			}
		}
	}

	return &FetchRepositoriesResponse{
		Repositories: allRepos,
	}, nil
}

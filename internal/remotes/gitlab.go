package remotes

import (
	"fmt"

	"github.com/depscloud/indexer/internal/config"

	"github.com/sirupsen/logrus"

	"github.com/xanzy/go-gitlab"
)

// NewGitlabRemote constructs a new remote implementation that speaks with Gitlab
// for repository related information.
func NewGitlabRemote(cfg *config.Gitlab) (Remote, error) {
	options := make([]gitlab.ClientOptionFunc, 0)
	if baseURL := cfg.GetBaseUrl(); baseURL != "" {
		options = append(options, gitlab.WithBaseURL(baseURL))
	}

	var client *gitlab.Client
	var err error

	if private := cfg.GetPrivate(); private != nil {
		client, err = gitlab.NewClient(private.GetToken(), options...)
	} else {
		err = fmt.Errorf("no auth method provided")
	}

	if err != nil {
		return nil, err
	}

	return &gitlabRemote{
		config: cfg,
		client: client,
	}, nil
}

var _ Remote = &gitlabRemote{}

type gitlabRemote struct {
	config *config.Gitlab
	client *gitlab.Client
}

func (r *gitlabRemote) FetchRepositories(request *FetchRepositoriesRequest) (*FetchRepositoriesResponse, error) {
	cloneConfig := r.config.GetClone()

	// if clone config is nil, fall back
	if cloneConfig == nil {
		cloneConfig = &config.Clone{
			Strategy: r.config.GetStrategy(),
		}
	}

	repositories := make([]*Repository, 0)
	groups := make(map[string]bool, 0)
	for _, group := range r.config.GetGroups() {
		groups[group] = true
	}

	logrus.Infof("[remotes.gitlab] fetching groups")
	page := 1
	for page > 0 {
		grps, resp, err := r.client.Groups.ListGroups(&gitlab.ListGroupsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})

		if err != nil {
			logrus.Errorf("[remotes.gitlab] encountered err while fetching groups %v", err)
			break
		}

		for _, group := range grps {
			groups[group.Path] = true
		}

		page = resp.NextPage
	}

	for _, user := range r.config.Users {
		logrus.Infof("[remotes.gitlab] fetching projects for user: %s", user)

		page = 1
		for page > 0 {
			projects, resp, err := r.client.Projects.ListUserProjects(user, &gitlab.ListProjectsOptions{
				ListOptions: gitlab.ListOptions{
					Page:    page,
					PerPage: 100,
				},
			})

			if err != nil {
				logrus.Errorf("[remotes.gitlab] encountered err while fetching projects for user %s, %v", user, err)
				break
			}

			urls := make([]*Repository, len(projects))

			for i, project := range projects {
				if cloneConfig.GetStrategy() == config.CloneStrategy_HTTP {
					urls[i] = &Repository{
						RepositoryURL: project.HTTPURLToRepo,
						Clone:         cloneConfig,
					}
				} else {
					urls[i] = &Repository{
						RepositoryURL: project.SSHURLToRepo,
						Clone:         cloneConfig,
					}
				}
			}

			repositories = append(repositories, urls...)

			page = resp.NextPage
		}
	}

	for group := range groups {
		logrus.Infof("[remotes.gitlab] fetching projects for group: %s", group)

		page = 1
		for page > 0 {
			projects, resp, err := r.client.Groups.ListGroupProjects(group, &gitlab.ListGroupProjectsOptions{
				ListOptions: gitlab.ListOptions{
					Page:    page,
					PerPage: 100,
				},
			})

			if err != nil {
				logrus.Errorf("[remotes.gitlab] encountered err while fetching projects for group %s, %v", group, err)
				break
			}

			urls := make([]*Repository, len(projects))

			for i, project := range projects {
				if cloneConfig.GetStrategy() == config.CloneStrategy_HTTP {
					urls[i] = &Repository{
						RepositoryURL: project.HTTPURLToRepo,
						Clone:         cloneConfig,
					}
				} else {
					urls[i] = &Repository{
						RepositoryURL: project.SSHURLToRepo,
						Clone:         cloneConfig,
					}
				}
			}

			repositories = append(repositories, urls...)

			page = resp.NextPage
		}
	}

	return &FetchRepositoriesResponse{
		Repositories: repositories,
	}, nil
}

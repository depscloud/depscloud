package remotes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/deps-cloud/indexer/internal/config"

	"github.com/nytlabs/gojee"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

// NewGenericRemote constructs a new remote endpoint that
func NewGenericRemote(cfg *config.Generic) Remote {
	return &genericRemote{
		config: cfg,
	}
}

var _ Remote = &genericRemote{}

type genericRemote struct {
	config *config.Generic
}

func (r *genericRemote) FetchRepositories(request *FetchRepositoriesRequest) (*FetchRepositoriesResponse, error) {
	cloneConfig := r.config.GetClone()

	tokens, err := jee.Lexer(r.config.Selector)
	if err != nil {
		return nil, err
	}

	parser, err := jee.Parser(tokens)
	if err != nil {
		return nil, err
	}

	logrus.Infof("[remotes.generic] requesting data from generic endpoint: %s", r.config.BaseUrl)

	repositories := make([]*Repository, 0)
	for page := 1; true; page++ {
		fullURL := fmt.Sprintf(
			"%s%s?%s=%d&%s=%d",
			r.config.BaseUrl,
			r.config.Path,
			r.config.PageParameter,
			page,
			r.config.PerPageParameter,
			r.config.PageSize,
		)

		resp, err := http.Get(fullURL)
		if err != nil {
			return nil, errors.Wrap(err,
				fmt.Sprintf("failed to get url: %s", fullURL))
		}

		if resp.StatusCode == http.StatusNotFound {
			logrus.Infof("[remotes.generic] encountered a 404. assuming end of data")
			break
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read body")
		}

		var umsg jee.BMsg
		if err := json.Unmarshal(body, &umsg); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal JSON")
		}

		result, err := jee.Eval(parser, umsg)
		if err != nil {
			return nil, errors.Wrapf(err,
				fmt.Sprintf("failed to extract response from page using selector: %s", r.config.Selector))
		}

		resultArray := result.([]interface{})
		for _, entry := range resultArray {
			entryString := entry.(string)
			repositories = append(repositories, &Repository{
				RepositoryURL: entryString,
				Clone:         cloneConfig,
			})
		}

		if int32(len(resultArray)) < r.config.PageSize {
			logrus.Infof("[remotes.generic] encountered an incomplete page. assuming end of data")
			break
		}
	}

	return &FetchRepositoriesResponse{
		Repositories: repositories,
	}, nil
}

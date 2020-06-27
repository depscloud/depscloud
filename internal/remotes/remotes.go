package remotes

import (
	"fmt"

	"github.com/depscloud/indexer/internal/config"
)

// ParseConfig is used to parse the account configuration and construct the
// necessary remote endpoint based on the configuration object.
func ParseConfig(configuration *config.Configuration) (Remote, error) {
	remotes := make([]Remote, len(configuration.Accounts))

	for i, account := range configuration.Accounts {
		var remote Remote
		var err error

		if generic := account.GetGeneric(); generic != nil {
			remote = NewGenericRemote(generic)
		} else if bitbucket := account.GetBitbucket(); bitbucket != nil {
			remote, err = NewBitbucketRemote(bitbucket)
		} else if github := account.GetGithub(); github != nil {
			remote, err = NewGithubRemote(github)
		} else if gitlab := account.GetGitlab(); gitlab != nil {
			remote, err = NewGitlabRemote(gitlab)
		} else if static := account.GetStatic(); static != nil {
			remote = NewStaticRemote(static)
		} else {
			err = fmt.Errorf("unrecognized account")
		}

		if err != nil {
			return nil, err
		}

		remotes[i] = remote
	}

	return NewCompositeRemote(remotes...), nil
}

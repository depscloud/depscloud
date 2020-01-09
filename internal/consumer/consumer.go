package consumer

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/deps-cloud/api/v1alpha/extractor"
	"github.com/deps-cloud/api/v1alpha/schema"
	"github.com/deps-cloud/api/v1alpha/tracker"
	"github.com/deps-cloud/discovery/pkg/remotes"

	"github.com/sirupsen/logrus"

	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

// RepositoryConsumer represent the contract for consuming repositories
type RepositoryConsumer interface {
	Consume(repository *remotes.Repository)
}

// NewConsumer creates a consumer process that is agnostic to the ingress channel.
func NewConsumer(
	authMethod transport.AuthMethod,
	desClient extractor.DependencyExtractorClient,
	sourceService tracker.SourceServiceClient,
) RepositoryConsumer {
	return &consumer{
		authMethod:    authMethod,
		desClient:     desClient,
		sourceService: sourceService,
	}
}

type consumer struct {
	authMethod    transport.AuthMethod
	desClient     extractor.DependencyExtractorClient
	sourceService tracker.SourceServiceClient
}

var _ RepositoryConsumer = &consumer{}

func (c *consumer) Consume(repository *remotes.Repository) {
	repourl := repository.RepositoryURL

	dir, err := ioutil.TempDir(os.TempDir(), "dis")
	if err != nil {
		logrus.Errorf("failed to create tempdir")
		return
	}

	fs := osfs.New(dir)
	gitfs, err := fs.Chroot(git.GitDirName)
	if err != nil {
		logrus.Errorf("failed to chroot for .git: %v", err)
		return
	}

	storage := filesystem.NewStorage(gitfs, cache.NewObjectLRUDefault())
	options := &git.CloneOptions{
		URL:   repourl,
		Depth: 1,
	}

	if repository.Clone != nil {
		if basic := repository.Clone.GetBasic(); basic != nil {
			auth := &http.BasicAuth{
				Username: basic.GetUsername(),
			}

			if password := basic.GetPassword(); password != nil {
				auth.Password = password.Value
			}

			options.Auth = auth
		} else if publicKey := repository.Clone.GetPublicKey(); publicKey != nil {
			user := publicKey.GetUser()
			if user == "" {
				user = "git"
			}

			password := publicKey.GetPassword()

			var keys *ssh.PublicKeys
			if privateKeyPath := publicKey.GetPrivateKeyPath(); privateKeyPath != nil {
				keys, err = ssh.NewPublicKeysFromFile(user, privateKeyPath.Value, password)
			} else if privateKey := publicKey.GetPrivateKey(); privateKey != nil {
				keys, err = ssh.NewPublicKeys(user, []byte(privateKey.GetValue()), password)
			}

			if keys == nil || err != nil {
				logrus.Errorf("[%s] failed to get public keys for repository", repourl)
				return
			}

			options.Auth = keys
		}
	} else if c.authMethod != nil {
		options.Auth = c.authMethod
	}

	logrus.Infof("[%s] cloning repository", repourl)
	_, err = git.Clone(storage, fs, options)

	if err != nil {
		logrus.Errorf("failed to clone: %v", err)
		return
	}

	logrus.Infof("[%s] walking file system", repourl)
	queue := []string{""}
	paths := make([]string, 0)

	for len(queue) > 0 {
		newQueue := make([]string, 0)
		size := len(queue)

		for i := 0; i < size; i++ {
			path := queue[i]

			finfos, err := fs.ReadDir(path)
			if err != nil {
				logrus.Errorf("failed to stat path: %v", err)
			}

			for _, finfo := range finfos {
				fpath := fs.Join(path, finfo.Name())
				if finfo.IsDir() {
					newQueue = append(newQueue, fpath)
				} else {
					paths = append(paths, fpath)
				}
			}
		}

		queue = newQueue
	}

	logrus.Infof("[%s] matching dependency files", repourl)
	matchedResponse, err := c.desClient.Match(context.Background(), &extractor.MatchRequest{
		Separator: string(filepath.Separator),
		Paths:     paths,
	})

	fileContents := make(map[string]string)
	for _, matched := range matchedResponse.MatchedPaths {
		file, err := fs.Open(matched)
		if err != nil {
			logrus.Warnf("failed to open file %s: %v", matched, err)
			continue
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			logrus.Warnf("failed to read file %s: %v", matched, err)
			continue
		}

		fileContents[matched] = string(data)
	}

	logrus.Infof("[%s] extracting dependencies", repourl)
	extractResponse, err := c.desClient.Extract(context.Background(), &extractor.ExtractRequest{
		Url:          repourl,
		Separator:    string(filepath.Separator),
		FileContents: fileContents,
	})

	if err != nil {
		logrus.Errorf("failed to extract deps from repo: %s", repourl)
		return
	}

	logrus.Infof("[%s] storing dependencies", repourl)
	_, err = c.sourceService.Track(context.Background(), &tracker.SourceRequest{
		Source: &schema.Source{
			Url: repourl,
		},
		ManagementFiles: extractResponse.GetManagementFiles(),
	})

	if err != nil {
		logrus.Errorf("failed to update deps for repo: %s, %v", repourl, err)
		return
	}

	logrus.Infof("[%s] cleaning up file system", repourl)
	if err := os.RemoveAll(dir); err != nil {
		logrus.Errorf("failed to cleanup scratch directory: %s", err.Error())
	}
}

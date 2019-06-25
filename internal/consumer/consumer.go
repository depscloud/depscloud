package consumer

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	desapi "github.com/deps-cloud/des/api"
	dtsapi "github.com/deps-cloud/dts/api"
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

// RepositoryConsumer represent the contract for consuming repositories
type RepositoryConsumer interface {
	Consume(repository string)
}

// NewConsumer creates a consumer process that is agnostic to the ingress channel.
func NewConsumer(
	authMethod transport.AuthMethod,
	desClient  desapi.DependencyExtractorClient,
	dtsClient  dtsapi.DependencyTrackerClient,
) RepositoryConsumer {
	return &consumer{
		authMethod: authMethod,
		desClient: desClient,
		dtsClient: dtsClient,
	}
}

type consumer struct {
	authMethod transport.AuthMethod
	desClient  desapi.DependencyExtractorClient
	dtsClient  dtsapi.DependencyTrackerClient
}

var _ RepositoryConsumer = &consumer{}

func (c *consumer) Consume(repository string) {
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
		URL:   repository,
		Depth: 1,
	}

	if c.authMethod != nil {
		options.Auth = c.authMethod
	}

	logrus.Infof("[%s] cloning repository", repository)
	_, err = git.Clone(storage, fs, options)

	if err != nil {
		logrus.Errorf("failed to clone: %v", err)
		return
	}

	logrus.Infof("[%s] walking file system", repository)
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

	logrus.Infof("[%s] matching dependency files", repository)
	matchedResponse, err := c.desClient.Match(context.Background(), &desapi.MatchRequest{
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

	logrus.Infof("[%s] extracting dependencies", repository)
	extractResponse, err := c.desClient.Extract(context.Background(), &desapi.ExtractRequest{
		Separator:    string(filepath.Separator),
		FileContents: fileContents,
	})

	if err != nil {
		logrus.Errorf("failed to extract deps from repo: %s", repository)
		return
	}

	logrus.Infof("[%s] storing dependencies", repository)
	resp, err := c.dtsClient.Put(context.Background(), &dtsapi.PutRequest{
		SourceInformation: &dtsapi.SourceInformation{
			Url: repository,
		},
		ManagementFiles: extractResponse.ManagementFiles,
	})

	if err != nil {
		logrus.Errorf("failed to update deps for repo: %s, %v", repository, err)
		return
	}

	if resp.Code != http.StatusOK {
		logrus.Errorf("[%s] %s", repository, resp.Message)
	} else {
		logrus.Infof("[%s] %s", repository, resp.Message)
	}

	logrus.Infof("[%s] cleaning up file system", repository)
	if err := os.RemoveAll(dir); err != nil {
		logrus.Errorf("failed to cleanup scratch directory: %s", err.Error())
	}
}

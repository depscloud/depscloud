package consumer

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/depscloud/api/v1alpha/extractor"
	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/depscloud/indexer/internal/remotes"
	"github.com/depscloud/depscloud/internal/logger"

	"go.uber.org/zap"

	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

// RepositoryConsumer represent the contract for consuming repositories
type RepositoryConsumer interface {
	Consume(ctx context.Context, repository *remotes.Repository)
}

// NewConsumer creates a consumer process that is agnostic to the ingress channel.
func NewConsumer(
	extractorService extractor.DependencyExtractorClient,
	sourceService tracker.SourceServiceClient,
) RepositoryConsumer {
	return &consumer{
		extractorService: extractorService,
		sourceService:    sourceService,
	}
}

type consumer struct {
	extractorService extractor.DependencyExtractorClient
	sourceService    tracker.SourceServiceClient
}

var _ RepositoryConsumer = &consumer{}

func (c *consumer) Consume(ctx context.Context, repository *remotes.Repository) {
	repoURL := repository.RepositoryURL

	log := logger.Extract(ctx)
	log = log.With(zap.String("repoURL", repoURL))

	dir, err := ioutil.TempDir(os.TempDir(), "dis")
	if err != nil {
		log.Error("failed to create temp directory",
			zap.String("dir", dir),
			zap.Error(err))
		return
	}

	// ensure proper cleanup
	defer func() {
		log.Info("cleaning up file system",
			zap.String("dir", dir))

		if err := os.RemoveAll(dir); err != nil {
			log.Error("failed to cleanup scratch directory",
				zap.String("dir", dir),
				zap.Error(err))
		}
	}()

	fs := osfs.New(dir)
	gitfs, err := fs.Chroot(git.GitDirName)
	if err != nil {
		log.Error("failed to chroot for .git",
			zap.String("dir", dir),
			zap.Error(err))
		return
	}

	storage := filesystem.NewStorage(gitfs, cache.NewObjectLRUDefault())
	options := &git.CloneOptions{
		URL:   repoURL,
		Depth: 1,
	}

	if repository.Clone != nil {
		if basic := repository.Clone.GetBasic(); basic != nil {
			options.Auth = &http.BasicAuth{
				Username: basic.GetUsername(),
				Password: basic.GetPassword(),
			}
		} else if publicKey := repository.Clone.GetPublicKey(); publicKey != nil {
			user := publicKey.GetUser()
			if user == "" {
				user = "git"
			}

			password := publicKey.GetPassword()

			var keys *ssh.PublicKeys
			if privateKeyPath := publicKey.GetPrivateKeyPath(); privateKeyPath != "" {
				keys, err = ssh.NewPublicKeysFromFile(user, privateKeyPath, password)
			} else if privateKey := publicKey.GetPrivateKey(); privateKey != "" {
				keys, err = ssh.NewPublicKeys(user, []byte(privateKey), password)
			}

			if keys == nil || err != nil {
				log.Error("failed to get public keys for repository",
					zap.Error(err))
				return
			}

			options.Auth = keys
		}
	}

	log.Info("cloning repository")
	repo, err := git.Clone(storage, fs, options)

	if err != nil {
		log.Error("failed to clone repository", zap.Error(err))
		return
	}

	log.Info("walking file system")
	queue := []string{""}
	paths := make([]string, 0)

	for len(queue) > 0 {
		newQueue := make([]string, 0)
		size := len(queue)

		for i := 0; i < size; i++ {
			path := queue[i]

			finfos, err := fs.ReadDir(path)
			if err != nil {
				log.Error("failed to stat path", zap.Error(err))
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

	log.Info("matching dependency files")
	matchedResponse, err := c.extractorService.Match(context.Background(), &extractor.MatchRequest{
		Separator: string(filepath.Separator),
		Paths:     paths,
	})

	if err != nil {
		log.Error("failed to match paths for repository")
		return
	}

	fileContents := make(map[string]string)
	for _, matched := range matchedResponse.MatchedPaths {
		file, err := fs.Open(matched)
		if err != nil {
			log.Warn("failed to open file",
				zap.String("matched", matched),
				zap.Error(err))
			continue
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Warn("failed to read file",
				zap.String("matched", matched),
				zap.Error(err))
			continue
		}

		fileContents[matched] = string(data)
	}

	log.Info("extracting dependencies")
	extractResponse, err := c.extractorService.Extract(context.Background(), &extractor.ExtractRequest{
		Url:          repoURL,
		Separator:    string(filepath.Separator),
		FileContents: fileContents,
	})

	if err != nil {
		log.Error("failed to extract deps from repo", zap.Error(err))
		return
	}

	ref := ""
	if head, err := repo.Head(); err == nil {
		ref = head.Name().String()
	}

	log.Info("storing dependencies")
	_, err = c.sourceService.Track(context.Background(), &tracker.SourceRequest{
		Source: &schema.Source{
			Url:  repoURL,
			Kind: "repository",
			Ref:  ref,
		},
		ManagementFiles: extractResponse.GetManagementFiles(),
	})

	if err != nil {
		log.Error("failed to update deps for repo", zap.Error(err))
	}
}

package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	desapi "github.com/deps-cloud/des/api"
	dtsapi "github.com/deps-cloud/dts/api"
	rdsapi "github.com/deps-cloud/rds/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

func panicIff(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func dial(target string) *grpc.ClientConn {
	dialOptions := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}

	cc, err := grpc.Dial(target, dialOptions...)
	panicIff(err)

	return cc
}

// NewConsumer creates a consumer process that is agnostic to the ingress channel.
func NewConsumer(
	authMethod transport.AuthMethod,
	desClient desapi.DependencyExtractorClient,
	dtsClient dtsapi.DependencyTrackerClient,
) func(string) {
	return func(url string) {
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
			URL:   url,
			Depth: 1,
		}

		if authMethod != nil {
			options.Auth = authMethod
		}

		logrus.Infof("[%s] cloning repository", url)
		_, err = git.Clone(storage, fs, options)

		if err != nil {
			logrus.Errorf("failed to clone: %v", err)
			return
		}

		queue := []string{ "" }
		paths := make([]string, 0)

		for ; len(queue) > 0 ; {
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

		logrus.Infof("[%s] matching dependency files", url)
		matchedResponse, err := desClient.Match(context.Background(), &desapi.MatchRequest{
			Separator: string(filepath.Separator),
			Paths: paths,
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

		logrus.Infof("[%s] extracting dependencies", url)
		extractResponse, err := desClient.Extract(context.Background(), &desapi.ExtractRequest{
			Separator: string(filepath.Separator),
			FileContents: fileContents,
		})

		if err != nil {
			logrus.Errorf("failed to extract deps from repo: %s", url)
			return
		}

		logrus.Infof("[%s] storing dependencies", url)
		resp, err := dtsClient.Put(context.Background(), &dtsapi.PutRequest{
			SourceInformation: &dtsapi.SourceInformation{
				Url: url,
			},
			ManagementFiles: extractResponse.ManagementFiles,
		})

		if err != nil {
			logrus.Errorf("failed to update deps for repo: %s, %v", url, err)
			return
		}

		if resp.Code != http.StatusOK {
			logrus.Infof("[%s] %s", url, resp.Message)
		} else {
			logrus.Errorf("[%s] %s", url, resp.Message)
		}

		if err := os.RemoveAll(dir); err != nil {
			logrus.Errorf("failed to cleanup scratch directory: %s", err.Error())
		}
	}
}

// NewWorker encapsulates logic for pulling information off a channel and invoking the consumer
func NewWorker(in chan string, consumer func(string)) {
	for str := range in {
		consumer(str)
	}
}

func main() {
	workers := 5
	rdsAddress := "rds:8090"
	desAddress := "des:8090"
	dtsAddress := "dts:8090"

	sshUser := "git"
	sshKeyPath := ""

	cmd := &cobra.Command{
		Use: "dis",
		Short: "dependency indexing service",
		Run: func(cmd *cobra.Command, args []string) {
			rdsClient := rdsapi.NewRepositoryDiscoveryClient(dial(rdsAddress))
			desClient := desapi.NewDependencyExtractorClient(dial(desAddress))
			dtsClient := dtsapi.NewDependencyTrackerClient(dial(dtsAddress))

			var authMethod transport.AuthMethod

			if len(sshKeyPath) > 0 {
				logrus.Infof("[main] loading ssh key")
				var err error
				authMethod, err = ssh.NewPublicKeysFromFile(sshUser, sshKeyPath, "")
				panicIff(err)
			}

			repositories := make(chan string, workers)

			consumer := NewConsumer(authMethod, desClient, dtsClient)
			for i := 0; i < workers; i++ {
				go NewWorker(repositories, consumer)
			}

			for {
				listResponse, err := rdsClient.List(context.Background(), &rdsapi.ListRepositoriesRequest{})
				if err != nil {
					logrus.Errorf("[main] encountered an error trying to list repositories from rds: %v", err)

					time.Sleep(30 * time.Second)
				} else {
					for _, repository := range listResponse.Repositories {
						repositories <- repository
					}

					time.Sleep(1 * time.Hour)
				}
			}
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&workers, "workers", workers, "(optional) number of workers to process repositories")
	flags.StringVar(&rdsAddress, "rds-address", rdsAddress, "(optional) address to rds")
	flags.StringVar(&desAddress, "des-address", desAddress, "(optional) address to des")
	flags.StringVar(&dtsAddress, "dts-address", dtsAddress, "(optional) address to dts")
	flags.StringVar(&sshUser, "ssh-user", sshUser, "(optional) the ssh user, typically git")
	flags.StringVar(&sshKeyPath, "ssh-keypath", sshKeyPath, "(optional) the path to the ssh key file")

	if err := cmd.Execute(); err != nil {
		panic(err.Error())
	}
}

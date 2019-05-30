package main

import (
	"context"
	desapi "github.com/deps-cloud/des/api"
	dtsapi "github.com/deps-cloud/dts/api"
	rdsapi "github.com/deps-cloud/rds/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"io/ioutil"
	"path/filepath"
	"time"
)

func dial(target string) *grpc.ClientConn {
	dialOptions := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}

	cc, err := grpc.Dial(target, dialOptions...)
	if err != nil {
		panic(err)
	}

	return cc
}

// NewConsumer creates a consumer process that is agnostic to the ingress channel.
func NewConsumer(
	desClient desapi.DependencyExtractorClient,
	dtsClient dtsapi.DependencyTrackingServiceClient,
) func(string) {
	return func(url string) {
		fs := memfs.New()
		gitfs, err := fs.Chroot(git.GitDirName)
		if err != nil {
			logrus.Errorf("failed to chroot for .git: %v", err)
			return
		}

		storage := filesystem.NewStorage(gitfs, cache.NewObjectLRUDefault())

		_, err = git.Clone(storage, fs, &git.CloneOptions{
			URL: 	url,
			Depth: 	1,
		})

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

		extractResponse, err := desClient.Extract(context.Background(), &desapi.ExtractRequest{
			Separator: string(filepath.Separator),
			FileContents: fileContents,
		})

		if err != nil {
			logrus.Errorf("failed to extract deps from repo: %s", url)
			return
		}

		_, err = dtsClient.Put(context.Background(), &dtsapi.PutRequest{
			SourceInformation: &dtsapi.SourceInformation{
				Url: url,
			},
			ManagementFiles: extractResponse.ManagementFiles,
		})

		if err != nil {
			logrus.Errorf("failed to update deps for repo: %s", url)
			return
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

	cmd := &cobra.Command{
		Use: "dis",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			rdsClient := rdsapi.NewRepositoryDiscoveryServiceClient(dial(rdsAddress))
			desClient := desapi.NewDependencyExtractorClient(dial(desAddress))
			dtsClient := dtsapi.NewDependencyTrackingServiceClient(dial(dtsAddress))

			repositories := make(chan string, workers)

			consumer := NewConsumer(desClient, dtsClient)
			for i := 0; i < workers; i++ {
				go NewWorker(repositories, consumer)
			}

			for {
				listResponse, err := rdsClient.List(context.Background(), &rdsapi.ListRepositoriesRequest{})
				if err != nil {
					logrus.Errorf("encountered an error trying to list repositories from rds: %v", err)

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

	if err := cmd.Execute(); err != nil {
		panic(err.Error())
	}
}

package main

import (
	"os"
	"sync"
	"time"

	"github.com/deps-cloud/discovery/pkg/config"
	"github.com/deps-cloud/discovery/pkg/remotes"
	desapi "github.com/deps-cloud/extractor/api"
	"github.com/deps-cloud/indexer/internal/consumer"
	"github.com/deps-cloud/tracker/api/v1alpha"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"google.golang.org/grpc"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
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

// NewWorker encapsulates logic for pulling information off a channel and invoking the consumer
func NewWorker(in chan string, done chan bool, rc consumer.RepositoryConsumer) {
	for str := range in {
		rc.Consume(str)
		done <- true
	}
}

// run is an internal method that represents a single pass over the set of repositories returned from the discovery service.
func run(remote remotes.Remote, repositories chan string, done chan bool) error {
	repos, err := remote.ListRepositories()
	if err != nil {
		return err
	}

	// wait for the done goroutine to finish
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(size int, wg *sync.WaitGroup) {
		for i := 0; i < size; i++ {
			<-done
		}
		wg.Done()
	}(len(repos), wg)

	for _, repository := range repos {
		repositories <- repository
	}

	wg.Wait()
	return nil
}

func main() {
	cron := false
	workers := 5
	rdsAddress := "rds:8090"
	desAddress := "des:8090"
	dtsAddress := "dts:8090"

	rdsConfigPath := ""

	sshUser := "git"
	sshKeyPath := ""

	cmd := &cobra.Command{
		Use:   "indexer",
		Short: "dependency indexing service",
		Run: func(cmd *cobra.Command, args []string) {
			desClient := desapi.NewDependencyExtractorClient(dial(desAddress))
			sourceService := v1alpha.NewSourceServiceClient(dial(dtsAddress))

			rdsConfig := &config.Configuration{
				Accounts: []*config.Account{
					{Rds: &config.Rds{Target: rdsAddress}},
				},
			}

			if len(rdsConfigPath) > 0 {
				var err error
				rdsConfig, err = config.Load(rdsConfigPath)
				panicIff(err)
			}

			remote, err := remotes.ParseConfig(rdsConfig)
			panicIff(err)

			var authMethod transport.AuthMethod

			if len(sshKeyPath) > 0 {
				logrus.Infof("[main] loading ssh key")
				var err error
				authMethod, err = ssh.NewPublicKeysFromFile(sshUser, sshKeyPath, "")
				panicIff(err)
			}

			repositories := make(chan string, workers)
			done := make(chan bool, workers)

			rc := consumer.NewConsumer(authMethod, desClient, sourceService)
			for i := 0; i < workers; i++ {
				go NewWorker(repositories, done, rc)
			}

			if cron {
				logrus.Infof("[main] running as cron")
				if err := run(remote, repositories, done); err != nil {
					logrus.Errorf("[main] encountered an error trying to list repositories from rds: %v", err)
					os.Exit(1)
				}
			} else {
				logrus.Infof("[main] running as daemon")
				for {
					sleep := time.Hour

					if err := run(remote, repositories, done); err != nil {
						logrus.Errorf("[main] encountered an error trying to list repositories from rds: %v", err)
						sleep = 30 * time.Second
					}

					time.Sleep(sleep)
				}
			}
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&cron, "cron", cron, "(optional) run the process as a cron job instead of a daemon")
	flags.IntVar(&workers, "workers", workers, "(optional) number of workers to process repositories")
	flags.StringVar(&rdsConfigPath, "rds-config", rdsConfigPath, "(optional) path to the rds config file")
	flags.StringVar(&rdsAddress, "rds-address", rdsAddress, "(optional) address to rds")
	flags.StringVar(&desAddress, "des-address", desAddress, "(optional) address to des")
	flags.StringVar(&dtsAddress, "dts-address", dtsAddress, "(optional) address to dts")
	flags.StringVar(&sshUser, "ssh-user", sshUser, "(optional) the ssh user, typically git")
	flags.StringVar(&sshKeyPath, "ssh-keypath", sshKeyPath, "(optional) the path to the ssh key file")

	if err := cmd.Execute(); err != nil {
		panic(err.Error())
	}
}

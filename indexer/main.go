package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/depscloud/api/v1alpha/extractor"
	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/depscloud/indexer/internal/config"
	"github.com/depscloud/depscloud/indexer/internal/consumer"
	"github.com/depscloud/depscloud/indexer/internal/remotes"
	"github.com/depscloud/depscloud/internal/client"
	"github.com/depscloud/depscloud/internal/v"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"

	_ "google.golang.org/grpc/health"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

// variables set during build using -X ldflag by goreleaser
var version string
var commit string
var date string

// NewWorker encapsulates logic for pulling information off a channel and invoking the consumer
func NewWorker(repositories chan *remotes.Repository, wg *sync.WaitGroup, rc consumer.RepositoryConsumer) {
	for repository := range repositories {
		rc.Consume(repository)
		wg.Done()
	}
}

type indexerConfig struct {
	workers    int
	configPath string
	sshUser    string
	sshKeyPath string
}

var description = strings.TrimSpace(`
   To learn more about how to configure the indexing layer, see our documentation.
   https://deps.cloud/docs/deploy/config/indexing/
`)

func main() {
	version := v.Info{Version: version, Commit: commit, Date: date}

	cfg := &indexerConfig{
		workers:    5,
		configPath: "",
		sshUser:    "git",
		sshKeyPath: "",
	}

	extractorConfig, extractorFlags := client.WithFlags("extractor", &client.Config{
		Address:       "extractor:8090",
		ServiceConfig: client.DefaultServiceConfig,
		LoadBalancer:  client.DefaultLoadBalancer,
		TLS:           false,
		TLSConfig:     &client.TLSConfig{},
	})

	trackerConfig, trackerFlags := client.WithFlags("tracker", &client.Config{
		Address:       "tracker:8090",
		ServiceConfig: client.DefaultServiceConfig,
		LoadBalancer:  client.DefaultLoadBalancer,
		TLS:           false,
		TLSConfig:     &client.TLSConfig{},
	})

	flags := []cli.Flag{
		&cli.IntFlag{
			Name:        "workers",
			Usage:       "number of workers to process repositories",
			Value:       cfg.workers,
			Destination: &cfg.workers,
			EnvVars:     []string{"WORKERS"},
		},
		&cli.StringFlag{
			Name:        "config",
			Usage:       "path to the config file",
			Value:       cfg.configPath,
			Destination: &cfg.configPath,
			EnvVars:     []string{"CONFIG_PATH"},
		},
		&cli.StringFlag{
			Name:        "ssh-user",
			Usage:       "the ssh user, typically git",
			Value:       cfg.sshUser,
			Destination: &cfg.sshUser,
			EnvVars:     []string{"SSH_USER"},
		},
		&cli.StringFlag{
			Name:        "ssh-keypath",
			Usage:       "the path to the ssh key file",
			Value:       cfg.sshKeyPath,
			Destination: &cfg.sshKeyPath,
			EnvVars:     []string{"SSH_KEYPATH"},
		},
	}
	flags = append(flags, extractorFlags...)
	flags = append(flags, trackerFlags...)

	app := &cli.App{
		Name:        "indexer",
		Usage:       "crawl sources and store extracted content",
		Description: description,
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Output version information",
				Action: func(c *cli.Context) error {
					versionString := fmt.Sprintf("%s %s", c.Command.Name, version)
					fmt.Println(versionString)
					return nil

				},
			},
		},
		Flags: flags,
		Action: func(context *cli.Context) error {
			extractorConn, err := client.Connect(extractorConfig)
			if err != nil {
				return err
			}
			defer extractorConn.Close()

			trackerConn, err := client.Connect(trackerConfig)
			if err != nil {
				return err
			}
			defer trackerConn.Close()

			extractorClient := extractor.NewDependencyExtractorClient(extractorConn)
			sourceService := tracker.NewSourceServiceClient(trackerConn)

			var remoteConfig *config.Configuration

			if len(cfg.configPath) > 0 {
				remoteConfig, err = config.Load(cfg.configPath)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("--config must be provided")
			}

			remote, err := remotes.ParseConfig(remoteConfig)
			if err != nil {
				return err
			}

			var authMethod transport.AuthMethod

			if len(cfg.sshKeyPath) > 0 {
				logrus.Infof("[main] loading ssh key")
				authMethod, err = ssh.NewPublicKeysFromFile(cfg.sshUser, cfg.sshKeyPath, "")
				if err != nil {
					return err
				}
			}

			resp, err := remote.FetchRepositories(&remotes.FetchRepositoriesRequest{})
			if err != nil {
				return err
			}

			// start a wait group to track remaining work
			wg := &sync.WaitGroup{}
			wg.Add(len(resp.Repositories))

			repositories := make(chan *remotes.Repository, cfg.workers)
			defer close(repositories)

			rc := consumer.NewConsumer(authMethod, extractorClient, sourceService)
			for i := 0; i < cfg.workers; i++ {
				go NewWorker(repositories, wg, rc)
			}

			// feed until there are no more left
			for _, repository := range resp.Repositories {
				repositories <- repository
			}

			// wait for all work to be done
			wg.Wait()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

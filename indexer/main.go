package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/depscloud/api/v1beta"
	"github.com/depscloud/depscloud/indexer/internal/config"
	"github.com/depscloud/depscloud/indexer/internal/consumer"
	"github.com/depscloud/depscloud/indexer/internal/remotes"
	"github.com/depscloud/depscloud/internal/client"
	"github.com/depscloud/depscloud/internal/eventlp"
	"github.com/depscloud/depscloud/internal/logger"
	"github.com/depscloud/depscloud/internal/v"

	"github.com/urfave/cli/v2"

	_ "google.golang.org/grpc/health"
)

// variables set during build using -X ldflag by goreleaser
var version string
var commit string
var date string

type indexerConfig struct {
	workers    int
	configPath string
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
	}

	loggerConfig, loggerFlags := logger.WithFlags(logger.DefaultConfig())

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
		// TODO: remove these. right now, they're still referenced in the latest helm chart
		&cli.StringFlag{
			Name: "ssh-keypath",
		},
		&cli.StringFlag{
			Name: "ssh-user",
		},
	}

	flags = append(flags, loggerFlags...)
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
		Action: func(c *cli.Context) error {
			if len(cfg.configPath) == 0 {
				return fmt.Errorf("--config must be provided")
			}

			remoteConfig, err := config.Load(cfg.configPath)
			if err != nil {
				return err
			}

			log := logger.MustGetLogger(loggerConfig)
			ctx := logger.ToContext(c.Context, log)

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

			extractionService := v1beta.NewManifestExtractionServiceClient(extractorConn)
			storageService := v1beta.NewManifestStorageServiceClient(trackerConn)
			rc := consumer.NewConsumer(extractionService, storageService)

			remote, err := remotes.ParseConfig(remoteConfig)
			if err != nil {
				return err
			}

			resp, err := remote.FetchRepositories(&remotes.FetchRepositoriesRequest{
				Context: ctx,
			})
			if err != nil {
				return err
			}

			eventLoop := eventlp.New()
			for i := 0; i < cfg.workers; i++ {
				go eventLoop.Start(ctx)
			}

			for _, repository := range resp.Repositories {
				repo := repository
				_ = eventLoop.Submit(func(ctx context.Context) {
					rc.Consume(ctx, repo)
				})
			}

			// wait for all work to be done
			return eventLoop.GracefullyStop()
		},
	}

	_ = app.Run(os.Args)
}

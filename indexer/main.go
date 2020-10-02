package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/depscloud/api/v1alpha/extractor"
	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/depscloud/indexer/internal/config"
	"github.com/depscloud/depscloud/indexer/internal/consumer"
	"github.com/depscloud/depscloud/indexer/internal/remotes"
	"github.com/depscloud/depscloud/internal/mux"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/health"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

// variables set during build using -X ldflag by goreleaser
var version string
var commit string
var date string

// https://github.com/grpc/grpc/blob/master/doc/service_config.md
const serviceConfigTemplate = `{
	"loadBalancingPolicy": "%s",
	"healthCheckConfig": {
		"serviceName": ""
	}
}`

func dial(target, certFile, keyFile, caFile, lbPolicy string) (*grpc.ClientConn, error) {
	serviceConfig := fmt.Sprintf(serviceConfigTemplate, lbPolicy)

	dialOptions := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(serviceConfig),
	}

	if len(certFile) > 0 {
		certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}

		certPool := x509.NewCertPool()
		bs, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, err
		}

		ok := certPool.AppendCertsFromPEM(bs)
		if !ok {
			return nil, fmt.Errorf("failed to append certs")
		}

		transportCreds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{certificate},
			RootCAs:      certPool,
		})

		dialOptions = append(dialOptions, grpc.WithTransportCredentials(transportCreds))
	} else {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	return grpc.Dial(target, dialOptions...)
}

func dialExtractor(cfg *indexerConfig) (*grpc.ClientConn, error) {
	return dial(cfg.extractorAddress,
		cfg.extractorCertPath, cfg.extractorKeyPath, cfg.extractorCAPath,
		cfg.extractorLBPolicy)
}

func dialTracker(cfg *indexerConfig) (*grpc.ClientConn, error) {
	return dial(cfg.trackerAddress,
		cfg.trackerCertPath, cfg.trackerKeyPath, cfg.trackerCAPath,
		cfg.trackerLBPolicy)
}

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

	extractorAddress  string
	extractorCertPath string
	extractorKeyPath  string
	extractorCAPath   string
	extractorLBPolicy string

	trackerAddress  string
	trackerCertPath string
	trackerKeyPath  string
	trackerCAPath   string
	trackerLBPolicy string

	sshUser    string
	sshKeyPath string
}

var description = strings.TrimSpace(`
   To learn more about how to configure the indexing layer, see our documentation.
   https://deps.cloud/docs/deploy/config/indexing/
`)

func main() {
	cfg := &indexerConfig{
		workers:           5,
		configPath:        "",
		extractorAddress:  "extractor:8090",
		extractorCertPath: "",
		extractorKeyPath:  "",
		extractorCAPath:   "",
		extractorLBPolicy: "round_robin",
		trackerAddress:    "tracker:8090",
		trackerCertPath:   "",
		trackerKeyPath:    "",
		trackerCAPath:     "",
		trackerLBPolicy:   "round_robin",
		sshUser:           "git",
		sshKeyPath:        "",
	}

	app := &cli.App{
		Name:        "indexer",
		Usage:       "crawl sources and store extracted content",
		Description: description,
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Output version information",
				Action: func(c *cli.Context) error {
					versionString := fmt.Sprintf("%s %s", c.Command.Name, mux.Version{Version: version, Commit: commit, Date: date})
					fmt.Println(versionString)
					return nil

				},
			},
		},
		Flags: []cli.Flag{
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
				Name:        "extractor-address",
				Usage:       "address to the extractor service",
				Value:       cfg.extractorAddress,
				Destination: &cfg.extractorAddress,
				EnvVars:     []string{"EXTRACTOR_ADDRESS"},
			},
			&cli.StringFlag{
				Name:        "extractor-cert",
				Usage:       "certificate used to enable TLS for the extractor",
				Value:       cfg.extractorCertPath,
				Destination: &cfg.extractorCertPath,
				EnvVars:     []string{"EXTRACTOR_CERT_PATH"},
			},
			&cli.StringFlag{
				Name:        "extractor-key",
				Usage:       "key used to enable TLS for the extractor",
				Value:       cfg.extractorKeyPath,
				Destination: &cfg.extractorKeyPath,
				EnvVars:     []string{"EXTRACTOR_KEY_PATH"},
			},
			&cli.StringFlag{
				Name:        "extractor-ca",
				Usage:       "ca used to enable TLS for the extractor",
				Value:       cfg.extractorCAPath,
				Destination: &cfg.extractorCAPath,
				EnvVars:     []string{"EXTRACTOR_CA_PATH"},
			},
			&cli.StringFlag{
				Name:        "extractor-lb",
				Usage:       "the load balancer policy to use for the extractor",
				Value:       cfg.extractorLBPolicy,
				Destination: &cfg.extractorLBPolicy,
				EnvVars:     []string{"EXTRACTOR_LBPOLICY"},
			},
			&cli.StringFlag{
				Name:        "tracker-address",
				Usage:       "address to the tracker service",
				Value:       cfg.trackerAddress,
				Destination: &cfg.trackerAddress,
				EnvVars:     []string{"TRACKER_ADDRESS"},
			},
			&cli.StringFlag{
				Name:        "tracker-cert",
				Usage:       "certificate used to enable TLS for the tracker",
				Value:       cfg.trackerCertPath,
				Destination: &cfg.trackerCertPath,
				EnvVars:     []string{"TRACKER_CERT_PATH"},
			},
			&cli.StringFlag{
				Name:        "tracker-key",
				Usage:       "key used to enable TLS for the tracker",
				Value:       cfg.trackerKeyPath,
				Destination: &cfg.trackerKeyPath,
				EnvVars:     []string{"TRACKER_KEY_PATH"},
			},
			&cli.StringFlag{
				Name:        "tracker-ca",
				Usage:       "ca used to enable TLS for the tracker",
				Value:       cfg.trackerCAPath,
				Destination: &cfg.trackerCAPath,
				EnvVars:     []string{"TRACKER_CA_PATH"},
			},
			&cli.StringFlag{
				Name:        "tracker-lb",
				Usage:       "the load balancer policy to use for the tracker",
				Value:       cfg.trackerLBPolicy,
				Destination: &cfg.trackerLBPolicy,
				EnvVars:     []string{"TRACKER_LBPOLICY"},
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
		},
		Action: func(context *cli.Context) error {
			extractorConn, err := dialExtractor(cfg)
			if err != nil {
				return err
			}
			defer extractorConn.Close()

			trackerConn, err := dialTracker(cfg)
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

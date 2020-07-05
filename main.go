package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/depscloud/api/v1alpha/extractor"
	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/indexer/internal/config"
	"github.com/depscloud/indexer/internal/consumer"
	"github.com/depscloud/indexer/internal/remotes"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/health"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

func exitIff(err error) {
	if err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}
}

// https://github.com/grpc/grpc/blob/master/doc/service_config.md
const serviceConfigTemplate = `{
	"loadBalancingPolicy": "%s",
	"healthCheckConfig": {
		"serviceName": ""
	}
}`

func dial(target, certFile, keyFile, caFile, lbPolicy string) *grpc.ClientConn {
	serviceConfig := fmt.Sprintf(serviceConfigTemplate, lbPolicy)

	dialOptions := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(serviceConfig),
	}

	if len(certFile) > 0 {
		certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
		exitIff(err)

		certPool := x509.NewCertPool()
		bs, err := ioutil.ReadFile(caFile)
		exitIff(err)

		ok := certPool.AppendCertsFromPEM(bs)
		if !ok {
			exitIff(fmt.Errorf("failed to append certs"))
		}

		transportCreds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{certificate},
			RootCAs:      certPool,
		})

		dialOptions = append(dialOptions, grpc.WithTransportCredentials(transportCreds))
	} else {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	cc, err := grpc.Dial(target, dialOptions...)
	exitIff(err)

	return cc
}

// NewWorker encapsulates logic for pulling information off a channel and invoking the consumer
func NewWorker(repositories chan *remotes.Repository, done chan bool, rc consumer.RepositoryConsumer) {
	for repository := range repositories {
		rc.Consume(repository)
		done <- true
	}
}

func run(remote remotes.Remote, repositories chan *remotes.Repository, done chan bool) error {
	resp, err := remote.FetchRepositories(&remotes.FetchRepositoriesRequest{})
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
	}(len(resp.Repositories), wg)

	for _, repository := range resp.Repositories {
		repositories <- repository
	}

	wg.Wait()
	return nil
}

func main() {
	workers := 5
	configPath := ""

	extractorAddress := "extractor:8090"
	extractorCert := ""
	extractorKey := ""
	extractorCA := ""
	extractorLBPolicy := "round_robin"

	trackerAddress := "tracker:8090"
	trackerCert := ""
	trackerKey := ""
	trackerCA := ""
	trackerLBPolicy := "round_robin"

	sshUser := "git"
	sshKeyPath := ""

	cmd := &cobra.Command{
		Use:   "indexer",
		Short: "dependency indexing service",
		Run: func(cmd *cobra.Command, args []string) {
			desClient := extractor.NewDependencyExtractorClient(dial(extractorAddress, extractorCert, extractorKey, extractorCA, extractorLBPolicy))
			sourceService := tracker.NewSourceServiceClient(dial(trackerAddress, trackerCert, trackerKey, trackerCA, trackerLBPolicy))

			var rdsConfig *config.Configuration

			if len(configPath) > 0 {
				var err error
				rdsConfig, err = config.Load(configPath)
				exitIff(err)
			}

			remote, err := remotes.ParseConfig(rdsConfig)
			exitIff(err)

			var authMethod transport.AuthMethod

			if len(sshKeyPath) > 0 {
				logrus.Infof("[main] loading ssh key")
				var err error
				authMethod, err = ssh.NewPublicKeysFromFile(sshUser, sshKeyPath, "")
				exitIff(err)
			}

			repositories := make(chan *remotes.Repository, workers)
			done := make(chan bool, workers)

			rc := consumer.NewConsumer(authMethod, desClient, sourceService)
			for i := 0; i < workers; i++ {
				go NewWorker(repositories, done, rc)
			}

			logrus.Infof("[main] running indexer cron")
			if err := run(remote, repositories, done); err != nil {
				logrus.Errorf("[main] encountered an error trying to list repositories from rds: %v", err)
				os.Exit(1)
			}
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&workers, "workers", workers, "(optional) number of workers to process repositories")
	flags.StringVar(&configPath, "config", configPath, "(optional) path to the config file")
	flags.StringVar(&configPath, "rds-config", configPath, "(deprecated) path to the rds config file")

	flags.StringVar(&extractorAddress, "extractor-address", extractorAddress, "(optional) address to the extractor service")
	flags.StringVar(&extractorCert, "extractor-cert", extractorCert, "(optional) certificate used to enable TLS for the extractor")
	flags.StringVar(&extractorKey, "extractor-key", extractorKey, "(optional) key used to enable TLS for the extractor")
	flags.StringVar(&extractorCA, "extractor-ca", extractorCA, "(optional) ca used to enable TLS for the extractor")
	flags.StringVar(&extractorLBPolicy, "extractor-lb", extractorLBPolicy, "(optional) the load balancer policy to use for the extractor")

	flags.StringVar(&trackerAddress, "tracker-address", trackerAddress, "(optional) address to the tracker service")
	flags.StringVar(&trackerCert, "tracker-cert", trackerCert, "(optional) certificate used to enable TLS for the tracker")
	flags.StringVar(&trackerKey, "tracker-key", trackerKey, "(optional) key used to enable TLS for the tracker")
	flags.StringVar(&trackerCA, "tracker-ca", trackerCA, "(optional) ca used to enable TLS for the tracker")
	flags.StringVar(&trackerLBPolicy, "tracker-lb", trackerLBPolicy, "(optional) the load balancer policy to use for the tracker")

	flags.StringVar(&sshUser, "ssh-user", sshUser, "(optional) the ssh user, typically git")
	flags.StringVar(&sshKeyPath, "ssh-keypath", sshKeyPath, "(optional) the path to the ssh key file")

	if err := cmd.Execute(); err != nil {
		panic(err.Error())
	}
}

package test

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

const (
	LocalstackImage = "localstack/localstack:latest"
	LocalstackPort  = "4566"
)

type LocalstackContainer struct {
	Container testcontainers.Container
	Context   context.Context
	URI       string
}

func SetupLocalstack(ctx context.Context) *LocalstackContainer {
	exposedPort := fmt.Sprintf("%s/tcp", LocalstackPort)
	req := testcontainers.ContainerRequest{
		Image:        LocalstackImage,
		ExposedPorts: []string{exposedPort},
		WaitingFor: wait.ForAll(
			wait.ForLog("Ready.").WithStartupTimeout(20*time.Second),
			wait.ForListeningPort(LocalstackPort),
		).WithDeadline(1 * time.Minute),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start localstack container")
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get container host for localstack")
	}

	mappedPort, err := container.MappedPort(ctx, LocalstackPort)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get mapped port for localstack")
	}

	uri := fmt.Sprintf("http://%s:%s", hostIP, mappedPort.Port())
	log.Info().Msgf("Localstack connection string: %s", uri)

	return &LocalstackContainer{
		Container: container,
		Context:   ctx,
		URI:       uri,
	}
}

func (l *LocalstackContainer) TeardownLocalstack() {
	if err := l.Container.Terminate(l.Context); err != nil {
		log.Error().Err(err).Msg("cannot close localstack")
	}
}

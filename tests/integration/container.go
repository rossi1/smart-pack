package integration

import (
	"context"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PsqlContainer struct {
	container *postgres.PostgresContainer
	ConnStr   string
}

func createPsqlContainer(ctx context.Context) (*PsqlContainer, error) {
	container, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase(databaseName),
		postgres.WithUsername(databaseUser),
		postgres.WithPassword(databasePassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := container.ConnectionString(ctx, databaseSSLMode)
	if err != nil {
		// Clean up container if we can't get connection string
		_ = container.Terminate(ctx)
		return nil, err
	}

	return &PsqlContainer{
		container: container,
		ConnStr:   connStr,
	}, nil
}

func (p *PsqlContainer) GetConnectionString(ctx context.Context) (string, error) {
	return p.container.ConnectionString(ctx, databaseSSLMode)
}

func (p *PsqlContainer) Terminate(ctx context.Context) error {
	return p.container.Terminate(ctx)
}

package testenv

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rossi1/smart-pack/adapters"
	"github.com/rossi1/smart-pack/cmd"
	cfg "github.com/rossi1/smart-pack/config"
)

func newTestDependencies(
	ctx context.Context,
	config *Config,
) (*cmd.Dependencies, error) {
	postgresClient, err := createPostgresClient(ctx, config.Cfg)
	if err != nil {
		return nil, err
	}

	manager := adapters.NewMigrationManager(config.Cfg.DatabaseURL,
		config.Cfg.DatabaseMigrationPath)

	if err := manager.LoadMigrations(); err != nil {
		return nil, err
	}

	return &cmd.Dependencies{
		DB: postgresClient,
	}, nil
}

func createPostgresClient(ctx context.Context, config *cfg.AppConfig) (*pgx.Conn, error) {
	dbOptions, err := pgx.ParseConfig(config.DatabaseURL)
	if err != nil {
		return nil, err
	}
	client, err := pgx.ConnectConfig(ctx, dbOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}

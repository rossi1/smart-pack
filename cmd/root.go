package cmd

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	appConfig "github.com/rossi1/smart-pack/config"
	"github.com/rossi1/smart-pack/pkg/config"
	"github.com/spf13/cobra"
)

var (
	logLevel string
	cfg      = &appConfig.AppConfig{}
)

var rootCmd = &cobra.Command{
	Use:   "smart-pack",
	Short: "smart-pack cli",
	Long:  `A CLI tool for managing smart-pack service`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		loader := configLoader()
		if err := loader.LoadConfig(".", cfg); err != nil {
			checkErr(err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&logLevel,
		"log-level",
		"l",
		"info",
		"Logging level. Default: info",
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

type Dependencies struct {
	DB *pgx.Conn
}

func (d *Dependencies) Close(ctx context.Context) error {
	if d.DB != nil {
		return d.DB.Close(ctx)
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func createPostgresConnection(ctx context.Context, cfg *appConfig.AppConfig) (*pgx.Conn, error) {
	dbOptions, err := pgx.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	db, err := pgx.ConnectConfig(ctx, dbOptions)
	if err != nil {
		return nil, err
	}

	// test connection
	_, err = db.Exec(ctx, "SELECT 1")
	if err != nil {
		db.Close(ctx)
		return nil, err
	}

	return db, nil
}

func configLoader() appConfig.Loader {
	return &config.ViperConfig{}
}

package cmd

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rossi1/smart-pack/adapters"
	smartCalculator "github.com/rossi1/smart-pack/adapters/smart_calculator"
	"github.com/rossi1/smart-pack/app"
	"github.com/rossi1/smart-pack/app/command"
	"github.com/rossi1/smart-pack/app/query"
	appConfig "github.com/rossi1/smart-pack/config"
	"github.com/rossi1/smart-pack/pkg/server"
	"github.com/rossi1/smart-pack/ports"
	"github.com/rossi1/smart-pack/ports/rest"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	port string
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "start api server",
	Run: func(cmd *cobra.Command, args []string) {
		startHTTP(cmd.Context(), cfg)
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	apiCmd.Flags().StringVarP(
		&port,
		"port",
		"p",
		"8080",
		"HTTP Server port to listen on.\nExample: smart-pack api -p 8080",
	)
}

func startHTTP(ctx context.Context, cfg *appConfig.AppConfig) {
	deps := initializeDependencies(ctx, cfg)
	defer safelyCloseDependencies(ctx, deps)

	application := NewApplication(ctx, cfg, deps)

	startRestServer(ctx, cfg, application)
}

func initializeDependencies(ctx context.Context, cfg *appConfig.AppConfig) *Dependencies {
	deps := newDependencies(ctx, cfg)
	if deps == nil {
		logrus.WithContext(ctx).Fatal("Failed to initialize dependencies")
	}
	return deps
}

func newDependencies(rootCtx context.Context, cfg *appConfig.AppConfig) *Dependencies {
	pgDB, err := createPostgresConnection(rootCtx, cfg)
	if err != nil {
		logrus.WithContext(rootCtx).Fatal("Error while connecting to postgres client", err)
	}

	return &Dependencies{
		DB: pgDB,
	}
}

func safelyCloseDependencies(ctx context.Context, deps *Dependencies) {
	if err := deps.Close(ctx); err != nil {
		logrus.Error("Error closing dependencies", err)
	}
}

func startRestServer(ctx context.Context, cfg *appConfig.AppConfig, application *app.Application) {
	server.RunHTTPServer(
		ctx,
		cfg,
		port,
		"/api",
		rest.SwaggerPath,
		func(router chi.Router) http.Handler {
			return ports.HandlerFromMux(rest.NewHTTPServer(application), router)
		},
	)
}

func NewApplication(ctx context.Context, cfg *appConfig.AppConfig, deps *Dependencies) *app.Application {
	logrus.WithContext(ctx).
		WithField("config", cfg).
		Info("Creating application")

	smartPackRepo := adapters.NewSmartPackRepository(deps.DB)
	packCalculator := smartCalculator.NewPackCalculator()

	return &app.Application{
		ErrorReporter: nil,
		AppConfig:     cfg,
		Commands: &app.Commands{
			SetPackSizes: command.NewSetPackSizesHandler(smartPackRepo),
		},
		Queries: &app.Queries{
			GetPackSizes: query.NewGetPackSizesHandler(smartPackRepo),
		},
		PackCalculator: packCalculator,
	}
}

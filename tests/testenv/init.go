package testenv

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	adapters "github.com/rossi1/smart-pack/adapters"
	"github.com/rossi1/smart-pack/app"
	"github.com/rossi1/smart-pack/cmd"
	appConfig "github.com/rossi1/smart-pack/config"
	"github.com/rossi1/smart-pack/pkg/server"
	"github.com/rossi1/smart-pack/ports"
	"github.com/rossi1/smart-pack/ports/rest"
)

func New(
	config *Config,
) (*TestApp, error) {
	ctx := context.Background()
	deps, err := newTestDependencies(ctx, config)
	if err != nil {
		return nil, err
	}

	application := cmd.NewApplication(ctx, config.Cfg, deps)

	psqlRepo := adapters.NewSmartPackRepository(deps.DB)

	repos := NewRepositories(
		psqlRepo,
	)

	httpServer := startTestHTTP(config.Cfg, application)
	restClient, err := ports.NewClientWithResponses(httpServer.URL)
	if err != nil {
		return nil, err
	}

	return &TestApp{
		close: func() {
			httpServer.Close()
			deps.Close(ctx)
		},
		Config:       config,
		Dependencies: deps,
		RestClient:   restClient,
		Repos:        repos,
		Application:  application,
	}, nil
}

func startTestHTTP(
	appCfg *appConfig.AppConfig,
	application *app.Application,
) *httptest.Server {
	return httptest.NewServer(
		server.GetRootRouter(appCfg, "/", rest.SwaggerPath, func(router chi.Router) http.Handler {
			return ports.HandlerFromMux(rest.NewHTTPServer(application), router)
		}),
	)
}

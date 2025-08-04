package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ghodss/yaml"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	appConfig "github.com/rossi1/smart-pack/config"
	"github.com/sirupsen/logrus"
)

const (
	teardownTimeout = 10
)

// RunHTTPServer configures and starts the listening http server.
func RunHTTPServer(
	ctx context.Context,
	cfg *appConfig.AppConfig,
	port, basePath, swaggerPath string,
	createHandler func(router chi.Router) http.Handler) {
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, syscall.SIGINT, syscall.SIGTERM)

	startCtx, startCancel := context.WithCancel(ctx)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           GetRootRouter(cfg, basePath, swaggerPath, createHandler),
		WriteTimeout:      cfg.ApplicationAPITimeout,
		ReadTimeout:       cfg.ApplicationAPITimeout,
		ReadHeaderTimeout: cfg.ApplicationAPITimeout,
	}

	go func() {
		<-osChan
		logrus.WithContext(startCtx).Debug("Server is shutting down...")

		ctxWithTimeout, cancel := context.WithTimeout(ctx, teardownTimeout*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctxWithTimeout); err != nil {
			logrus.WithContext(startCtx).
				WithError(err).
				Fatal("Could not gracefully shutdown the httpServer")
		}

		startCancel()
		logrus.WithContext(ctxWithTimeout).Info("Server stopped")
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.WithContext(startCtx).WithError(err).Fatal("Server shutting down, error occurred")
	}
}

func GetRootRouter(
	cfg *appConfig.AppConfig,
	basePath, swaggerPath string,
	apiHandler func(router chi.Router) http.Handler) http.Handler {
	apiRouter := chi.NewRouter()
	setMiddlewares(cfg, apiRouter)

	rootRouter := chi.NewRouter()

	rootRouter.Get(refineSwaggerPath(swaggerPath), Swagger)

	if apiHandler != nil {
		router := chi.NewRouter()
		setMiddlewares(cfg, router)
		rootRouter.Mount(basePath, apiHandler(router))
	}
	return rootRouter
}

func refineSwaggerPath(path string) string {
	s := path
	s = strings.TrimSuffix(s, ".json")
	s = strings.TrimSuffix(s, ".yaml")
	s = strings.TrimSuffix(s, ".yml")
	return fmt.Sprintf("%s{ext}", s)
}

func setMiddlewares(cfg *appConfig.AppConfig, router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(cfg.ApplicationAPITimeout))
	router.Use(middleware.DefaultLogger)

	addCorsMiddleware(router, cfg.CORSAllowedOrigins)

	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
		middleware.SetHeader("X-Frame-Options", "deny"),
	)
	router.Use(middleware.NoCache)
}

func addCorsMiddleware(router *chi.Mux, origins string) {
	allowedOrigins := strings.Split(origins, ";")
	if len(allowedOrigins) == 0 {
		return
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-App-ID"},
		ExposedHeaders:   []string{"Link", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsMiddleware.Handler)
}

func Swagger(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("./api/api.yml")
	if err != nil {
		logrus.WithContext(r.Context()).WithError(err).Info("failed to read api.yml")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	j, err := yaml.YAMLToJSON(file)

	if err != nil {
		logrus.WithContext(r.Context()).WithError(err).Info("failed to convert yaml to json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(j)
	if err != nil {
		logrus.WithContext(r.Context()).WithError(err).Info("failed to write api.yml to response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

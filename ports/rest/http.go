package rest

import (
	"net/http"

	"github.com/rossi1/smart-pack/app"
)

const SwaggerPath = "/api/v1/api-docs/swagger.json"

type HTTPServer struct {
	app *app.Application
}

func NewHTTPServer(application *app.Application) *HTTPServer {
	return &HTTPServer{
		app: application,
	}
}

func (h HTTPServer) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong")) //nolint:errcheck
}

func (h HTTPServer) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

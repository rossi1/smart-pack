package rest

import (
	"testing"
)

type testHTTPServer struct {
	api  HTTPServer
	deps *mockedDependencies
}

func newTestAPIServer(t *testing.T) testHTTPServer {
	deps := newMockedDeps(t)
	httpServer := NewHTTPServer(newTestApplication(deps))
	return testHTTPServer{
		api:  *httpServer,
		deps: deps,
	}
}

package testenv

import (
	"github.com/rossi1/smart-pack/adapters"
	"github.com/rossi1/smart-pack/app"
	"github.com/rossi1/smart-pack/cmd"
	restapi "github.com/rossi1/smart-pack/ports"
)

type Repositories struct {
	SmartPackRepository *adapters.SmartPackRepository
}

func NewRepositories(
	smartPackRepository *adapters.SmartPackRepository,
) *Repositories {
	return &Repositories{
		SmartPackRepository: smartPackRepository,
	}
}

type TestApp struct {
	Config       *Config
	Dependencies *cmd.Dependencies
	Application  *app.Application
	RestClient   *restapi.ClientWithResponses
	Repos        *Repositories
	close        func()
}

func (t *TestApp) Close() {
	t.close()
}

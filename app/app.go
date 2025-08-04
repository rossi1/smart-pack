package app

import (
	smartCalculator "github.com/rossi1/smart-pack/adapters/smart_calculator"
	"github.com/rossi1/smart-pack/app/command"
	"github.com/rossi1/smart-pack/app/query"
	appConfig "github.com/rossi1/smart-pack/config"
	"github.com/rossi1/smart-pack/domain"
)

type Application struct {
	ErrorReporter  domain.ErrorReporter
	AppConfig      *appConfig.AppConfig
	PackCalculator smartCalculator.PackCalculator
	Commands       *Commands
	Queries        *Queries
}

type Commands struct {
	SetPackSizes command.SetPackSizesHandler
}

type Queries struct {
	GetPackSizes query.GetPackSizesHandler
}

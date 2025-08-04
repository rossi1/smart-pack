package rest

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rossi1/smart-pack/adapters/smart_calculator"
	"github.com/rossi1/smart-pack/app"
	"github.com/rossi1/smart-pack/app/command"
	"github.com/rossi1/smart-pack/app/query"
)

type mockedDependencies struct {
	mockedSetPackSizesRepository command.SetPackSizesRepository
	mockedGetPackSizesRepository query.GetPackSizesRepository
	mockedPackCalculator         smart_calculator.PackCalculator
}

func newMockedDeps(t *testing.T) *mockedDependencies {
	ctrl := gomock.NewController(t)
	return &mockedDependencies{
		mockedSetPackSizesRepository: command.NewMockSetPackSizesRepository(ctrl),
		mockedGetPackSizesRepository: query.NewMockGetPackSizesRepository(ctrl),
		mockedPackCalculator:         smart_calculator.NewMockPackCalculator(ctrl),
	}
}

func newTestApplication(deps *mockedDependencies) *app.Application {
	return &app.Application{
		Commands: &app.Commands{
			SetPackSizes: command.NewSetPackSizesHandler(deps.mockedSetPackSizesRepository),
		},
		Queries: &app.Queries{
			GetPackSizes: query.NewGetPackSizesHandler(deps.mockedGetPackSizesRepository),
		},
		PackCalculator: deps.mockedPackCalculator,
	}
}

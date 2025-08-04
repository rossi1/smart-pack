package stories

import (
	"testing"

	"github.com/rossi1/smart-pack/tests/integration"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	integration.BaseSuite
}

func TestSmartPack(t *testing.T) {
	suite.Run(t, new(Suite))
}

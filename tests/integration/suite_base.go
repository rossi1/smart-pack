package integration

import (
	"context"
	"testing"
	"time"

	"github.com/rossi1/smart-pack/tests/testenv"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	timeout = 10 * time.Second
)

type BaseSuite struct {
	suite.Suite
	*testenv.TestApp
	psqlContainer *PsqlContainer
	ctx           context.Context
	cancel        context.CancelFunc
}

func (s *BaseSuite) Context() context.Context {
	return s.ctx
}

func (s *BaseSuite) SetupSuite() {
	if testing.Short() {
		s.T().Skip("Skipping integration tests")
	}

	s.T().Log("BaseSuite: SetupSuite Start")
	s.ctx, s.cancel = context.WithCancel(context.Background())

	r := require.New(s.T())
	var err error
	s.psqlContainer, err = createPsqlContainer(s.ctx)
	r.NoError(err)

	config, err := testenv.NewConfig(testenv.WithDatabaseURL(s.psqlContainer.ConnStr))
	r.NoError(err)

	s.TestApp, err = testenv.New(config)
	r.NoError(err)
	s.T().Log("BaseSuite: SetupSuite Done")
}

func (s *BaseSuite) TearDownSuite() {
	s.T().Log("BaseSuite: TearDownSuite Start")

	if s.TestApp != nil {
		s.Close()
	}
	if s.cancel != nil {
		s.cancel()
	}

	if s.psqlContainer != nil {
		ctx, cancel := context.WithTimeout(s.ctx, timeout)
		defer cancel()

		err := s.psqlContainer.Terminate(ctx)
		if err != nil {
			s.T().Logf("Failed to terminate container: %v", err)
		}
	}

	s.T().Log("BaseSuite: TearDownSuite Done")
}

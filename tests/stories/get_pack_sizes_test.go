package stories

import (
	"net/http"

	"github.com/stretchr/testify/require"
)

func (s *Suite) TestGetPackSizes() {
	r := require.New(s.T())

	resp, err := s.RestClient.GetPackSizesWithResponse(s.Context())
	r.NoError(err)
	r.Equal(http.StatusOK, resp.StatusCode())
	r.NotNil(resp.JSON200.PackSizes)
	r.GreaterOrEqual(len(resp.JSON200.PackSizes), 1)
}

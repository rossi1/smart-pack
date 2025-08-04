package stories

import (
	"net/http"

	restapi "github.com/rossi1/smart-pack/ports"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestSetPackSizes() {
	r := require.New(s.T())

	req := restapi.SetPackSizesRequest{
		PackSizes: []int{100, 200, 500},
	}

	resp, err := s.RestClient.SetPackSizesWithResponse(s.Context(), req)
	r.NoError(err)
	r.Equal(http.StatusOK, resp.StatusCode())
}

func (s *Suite) TestSetPackSizeError() {
	r := require.New(s.T())

	req := restapi.SetPackSizesRequest{
		PackSizes: []int{0, 100, 200, 500},
	}

	resp, err := s.RestClient.SetPackSizesWithResponse(s.Context(), req)
	r.NoError(err)
	r.Equal(http.StatusInternalServerError, resp.StatusCode())
}

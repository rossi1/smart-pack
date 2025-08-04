package stories

import (
	"net/http"
	"strconv"

	restapi "github.com/rossi1/smart-pack/ports"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestCalculatePack() {
	r := require.New(s.T())

	req := restapi.CalculateRequest{
		ItemsOrdered: 100,
	}

	resp, err := s.RestClient.CalculatePacksWithResponse(s.Context(), req)
	r.NoError(err)
	r.Equal(http.StatusOK, resp.StatusCode())
	r.NotNil(resp.JSON200)
	r.NotEmpty(resp.JSON200.Packs)
}

func (s *Suite) TestCalculateExactMatch() {
	r := require.New(s.T())

	req := restapi.CalculateRequest{ItemsOrdered: 750}
	resp, err := s.RestClient.CalculatePacksWithResponse(s.Context(), req)
	r.NoError(err)
	r.Equal(http.StatusOK, resp.StatusCode())

	expected := map[int]int{250: 1, 500: 1}

	actual := make(map[int]int)
	for k, v := range resp.JSON200.Packs {
		size, err := strconv.Atoi(k)
		r.NoError(err)
		actual[size] = v
	}

	r.Equal(expected, actual)
}

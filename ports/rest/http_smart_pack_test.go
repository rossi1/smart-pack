package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rossi1/smart-pack/adapters/smart_calculator"
	"github.com/rossi1/smart-pack/app/command"
	"github.com/rossi1/smart-pack/app/query"
	"github.com/rossi1/smart-pack/domain"
	"github.com/rossi1/smart-pack/ports"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
)

func TestCalculatePacks(t *testing.T) {
	testCases := []struct {
		Name         string
		RequestBody  ports.CalculateRequest
		MockFunc     func(server testHTTPServer)
		ResponseCode int
		ResponseBody *ports.PackSolution
	}{
		{
			Name: "invalid request (zero items_ordered)",
			RequestBody: ports.CalculateRequest{
				ItemsOrdered: 0,
			},
			ResponseCode: http.StatusBadRequest,
		},
		{
			Name: "internal error from GetPackSizes",
			RequestBody: ports.CalculateRequest{
				ItemsOrdered: 100,
			},
			MockFunc: func(server testHTTPServer) {
				server.deps.mockedGetPackSizesRepository.(*query.MockGetPackSizesRepository).
					EXPECT().
					GetPackSizes(gomock.Any()).
					Return(nil, errors.New("internal server error")).
					AnyTimes()
			},
			ResponseCode: http.StatusInternalServerError,
		},
		{
			Name: "success",
			RequestBody: ports.CalculateRequest{
				ItemsOrdered: 1200,
			},
			MockFunc: func(server testHTTPServer) {
				// Mock GetPackSizes to return pack sizes
				server.deps.mockedGetPackSizesRepository.(*query.MockGetPackSizesRepository).
					EXPECT().
					GetPackSizes(gomock.Any()).
					Return([]domain.SmartPack{
						{Size: 250}, {Size: 500}, {Size: 1000}, {Size: 2000}, {Size: 5000},
					}, nil).
					AnyTimes()

				// Mock PackCalculator.Calculate to return the expected PackSolution
				expectedSolution := &domain.PackSolution{
					ItemsOrdered: 1200,
					TotalItems:   1250,
					TotalPacks:   2,
					Packs:        map[int]int{1000: 1, 250: 1},
					PackDetails: []domain.PackDetail{
						{Size: 1000, Quantity: 1},
						{Size: 250, Quantity: 1},
					},
				}
				server.deps.mockedPackCalculator.(*smart_calculator.MockPackCalculator).
					EXPECT().
					Calculate(gomock.Any(), gomock.Any()).
					Return(expectedSolution, nil).
					AnyTimes()
			},
			ResponseCode: http.StatusOK,
			ResponseBody: &ports.PackSolution{
				ItemsOrdered: 1200,
				TotalItems:   1250,
				TotalPacks:   2,
				Packs:        map[string]int{"1000": 1, "250": 1},
				PackDetails: []ports.PackDetail{
					{Size: 1000, Quantity: 1},
					{Size: 250, Quantity: 1},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			testServer := newTestAPIServer(t)

			if tc.MockFunc != nil {
				tc.MockFunc(testServer)
			}

			data, err := json.Marshal(tc.RequestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewReader(data))
			req = req.WithContext(context.Background())
			req.Header.Set("Content-Type", "application/json")
			rw := httptest.NewRecorder()

			testServer.api.CalculatePacks(rw, req)

			if rw.Code != tc.ResponseCode {
				t.Logf("Response body: %s", rw.Body.String())
			}

			require.Equal(t, tc.ResponseCode, rw.Code)

			if tc.ResponseBody != nil {
				var actual ports.PackSolution
				err := json.Unmarshal(rw.Body.Bytes(), &actual)
				require.NoError(t, err)

				require.Equal(t, tc.ResponseBody.ItemsOrdered, actual.ItemsOrdered)
				require.Equal(t, tc.ResponseBody.TotalItems, actual.TotalItems)
				require.Equal(t, tc.ResponseBody.TotalPacks, actual.TotalPacks)
				require.Equal(t, tc.ResponseBody.Packs, actual.Packs)
				require.ElementsMatch(t, tc.ResponseBody.PackDetails, actual.PackDetails)
			}
		})
	}
}

func TestGetPackSizes(t *testing.T) {
	testCases := []struct {
		Name         string
		MockFunc     func(server testHTTPServer)
		ResponseCode int
		ResponseBody any
	}{

		{
			Name: "internal server error",
			MockFunc: func(server testHTTPServer) {
				server.deps.mockedGetPackSizesRepository.(*query.MockGetPackSizesRepository).
					EXPECT().GetPackSizes(gomock.Any()).
					Return(nil, errors.New("internal server error")).
					AnyTimes()
			},
			ResponseCode: http.StatusInternalServerError,
		},

		{
			Name: "success",
			MockFunc: func(server testHTTPServer) {
				server.deps.mockedGetPackSizesRepository.(*query.MockGetPackSizesRepository).
					EXPECT().GetPackSizes(gomock.Any()).
					Return([]domain.SmartPack{{Size: 250}, {Size: 500}, {Size: 1000}, {Size: 2000}, {Size: 5000}}, nil).
					AnyTimes()
			},
			ResponseCode: http.StatusOK,
			ResponseBody: ports.PackSizesResponse{
				PackSizes: []int{250, 500, 1000, 2000, 5000},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			testServer := newTestAPIServer(t)
			if tc.MockFunc != nil {
				tc.MockFunc(testServer)
			}

			r := httptest.NewRequest(http.MethodGet, "/api/pack-sizes", http.NoBody)
			r = r.WithContext(context.Background())
			r.Header.Set("content-type", "application/json")
			rw := httptest.NewRecorder()

			testServer.api.GetPackSizes(rw, r)

			if tc.ResponseCode == http.StatusOK {
				var resp ports.PackSizesResponse
				require.NoError(t, json.Unmarshal(rw.Body.Bytes(), &resp))
				require.Equal(t, tc.ResponseBody, resp)
			}
			require.Equal(t, tc.ResponseCode, rw.Code)
		})
	}
}

func TestSetPackSizes(t *testing.T) {
	testCases := []struct {
		Name         string
		MockFunc     func(server testHTTPServer)
		ResponseCode int
		RequestBody  ports.SetPackSizesRequest
	}{
		{
			Name: "internal server error",
			MockFunc: func(server testHTTPServer) {
				server.deps.mockedSetPackSizesRepository.(*command.MockSetPackSizesRepository).
					EXPECT().SetPackSizes(gomock.Any(), gomock.Any()).
					Return(errors.New("internal server error")).
					AnyTimes()
			},
			ResponseCode: http.StatusInternalServerError,
		},

		{
			Name: "success",
			MockFunc: func(server testHTTPServer) {
				server.deps.mockedSetPackSizesRepository.(*command.MockSetPackSizesRepository).
					EXPECT().SetPackSizes(gomock.Any(), gomock.Any()).
					Return(nil).
					AnyTimes()
			},
			ResponseCode: http.StatusOK,
			RequestBody: ports.SetPackSizesRequest{
				PackSizes: []int{250, 500, 1000, 2000, 5000},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			testServer := newTestAPIServer(t)
			if tc.MockFunc != nil {
				tc.MockFunc(testServer)
			}

			data, err := json.Marshal(&tc.RequestBody)
			require.NoError(t, err)

			r := httptest.NewRequest(http.MethodPost, "/api/pack-sizes", bytes.NewReader(data))
			r = r.WithContext(context.Background())
			r.Header.Set("content-type", "application/json")
			rw := httptest.NewRecorder()

			testServer.api.SetPackSizes(rw, r)

			require.Equal(t, tc.ResponseCode, rw.Code)
		})
	}
}

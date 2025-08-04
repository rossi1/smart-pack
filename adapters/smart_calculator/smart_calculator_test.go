package smart_calculator

import (
	"testing"

	"github.com/rossi1/smart-pack/domain"
	"github.com/stretchr/testify/require"
)

func TestPackCalculator_Calculate(t *testing.T) {
	calculator := NewPackCalculator()
	testCases := []struct {
		name       string
		order      int
		packSizes  []int
		expectErr  bool
		assertFunc func(t *testing.T, solution *domain.PackSolution)
	}{
		{
			name:      "invalid order zero",
			order:     0,
			packSizes: []int{250, 500},
			expectErr: true,
		},
		{
			name:      "empty pack sizes",
			order:     100,
			packSizes: []int{},
			expectErr: true,
		},
		{
			name:      "simple valid case",
			order:     1200,
			packSizes: []int{250, 500, 1000},
			expectErr: false,
			assertFunc: func(t *testing.T, solution *domain.PackSolution) {
				require.NotNil(t, solution)
				require.GreaterOrEqual(t, solution.TotalItems, 1200)
				require.Greater(t, solution.TotalPacks, 0)
				require.NotEmpty(t, solution.Packs)
				require.NotEmpty(t, solution.PackDetails)

				require.GreaterOrEqual(t, solution.TotalItems, solution.ItemsOrdered)
			},
		},
		{
			name:      "order less than smallest pack size",
			order:     200,
			packSizes: []int{250, 500},
			expectErr: false,
			assertFunc: func(t *testing.T, solution *domain.PackSolution) {
				require.NotNil(t, solution)
				require.GreaterOrEqual(t, solution.TotalItems, 200)
				require.Contains(t, solution.Packs, 250)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			solution, err := calculator.Calculate(tc.order, tc.packSizes)
			if tc.expectErr {
				require.Error(t, err)
				require.Nil(t, solution)
			} else {
				require.NoError(t, err)
				if tc.assertFunc != nil {
					tc.assertFunc(t, solution)
				}
			}
		})
	}
}

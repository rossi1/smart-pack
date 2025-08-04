package smart_calculator

import (
	"errors"
	"math"
	"sort"

	"github.com/rossi1/smart-pack/domain"
)

type PackCalculator interface {
	Calculate(order int, packSizes []int) (*domain.PackSolution, error)
}

type packCalculatorImpl struct{}

func NewPackCalculator() PackCalculator {
	return &packCalculatorImpl{}
}

func (c *packCalculatorImpl) Calculate(order int, packSizes []int) (*domain.PackSolution, error) {
	if order <= 0 {
		return nil, errors.New("order must be positive")
	}
	if len(packSizes) == 0 {
		return nil, errors.New("pack sizes empty")
	}

	// Sort pack sizes descending for better pruning and consistency
	sort.Sort(sort.Reverse(sort.IntSlice(packSizes)))

	result := findOptimalPacksMemo(order, packSizes)
	if len(result.Packs) == 0 {
		return nil, errors.New("cannot fulfill order with given pack sizes")
	}

	// Prepare PackDetails sorted descending by size
	sizes := make([]int, 0, len(result.Packs))
	for size := range result.Packs {
		sizes = append(sizes, size)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))

	details := make([]domain.PackDetail, 0, len(result.Packs))
	for _, size := range sizes {
		details = append(details, domain.PackDetail{
			Size:     size,
			Quantity: result.Packs[size],
		})
	}

	return &domain.PackSolution{
		ItemsOrdered: order,
		TotalItems:   result.TotalItems,
		TotalPacks:   result.TotalPacks,
		Packs:        result.Packs,
		PackDetails:  details,
	}, nil
}

type optimalPackSolution struct {
	Packs      map[int]int
	TotalItems int
	TotalPacks int
}

func findOptimalPacksMemo(order int, packSizes []int) optimalPackSolution {
	// DP approach: dp[i] = minimum items needed to fulfill at least i items
	maxCheck := order + packSizes[0] // Check up to order + largest pack
	dp := make([]int, maxCheck+1)
	parent := make([]int, maxCheck+1)

	// Initialize with impossible values
	for i := 1; i <= maxCheck; i++ {
		dp[i] = math.MaxInt32
		parent[i] = -1
	}
	dp[0] = 0

	for i := 0; i <= maxCheck; i++ {
		if dp[i] == math.MaxInt32 {
			continue
		}

		for _, packSize := range packSizes {
			next := i + packSize
			if next <= maxCheck && dp[next] > dp[i]+packSize {
				dp[next] = dp[i] + packSize
				parent[next] = packSize
			}
		}
	}

	bestItems := math.MaxInt32
	bestTarget := -1
	for i := order; i <= maxCheck; i++ {
		if dp[i] < bestItems {
			bestItems = dp[i]
			bestTarget = i
		}
	}

	if bestTarget == -1 {
		return optimalPackSolution{Packs: make(map[int]int)}
	}

	packs := make(map[int]int)
	current := bestTarget
	totalPacks := 0

	for current > 0 {
		packSize := parent[current]
		packs[packSize]++
		totalPacks++
		current -= packSize
	}

	return optimalPackSolution{
		Packs:      packs,
		TotalItems: bestItems,
		TotalPacks: totalPacks,
	}
}

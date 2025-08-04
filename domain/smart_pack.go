package domain

type SmartPack struct {
	Size int
}

type PackDetail struct {
	Size     int
	Quantity int
}

type PackSolution struct {
	ItemsOrdered int
	TotalItems   int
	TotalPacks   int
	Packs        map[int]int // size -> quantity
	PackDetails  []PackDetail
}

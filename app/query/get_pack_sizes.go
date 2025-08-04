package query

import (
	"context"

	"github.com/rossi1/smart-pack/common/decorator"
	"github.com/rossi1/smart-pack/domain"
)

type GetPackSizesQuery struct {
}

//go:generate mockgen -package=query -destination=get_pack_sizes.mock.go -source=get_pack_sizes.go
type GetPackSizesRepository interface {
	GetPackSizes(ctx context.Context) ([]domain.SmartPack, error)
}

type GetPackSizesHandler decorator.QueryHandler[*GetPackSizesQuery, []domain.SmartPack]

type getPackSizesHandler struct {
	repo GetPackSizesRepository
}

func NewGetPackSizesHandler(repo GetPackSizesRepository) GetPackSizesHandler {
	return decorator.ApplyQueryDecorators[*GetPackSizesQuery, []domain.SmartPack](&getPackSizesHandler{
		repo: repo,
	})
}

func (h *getPackSizesHandler) Handle(ctx context.Context, q *GetPackSizesQuery) ([]domain.SmartPack, error) {
	return h.repo.GetPackSizes(ctx)
}

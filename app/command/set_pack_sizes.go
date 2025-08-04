package command

import (
	"context"

	"github.com/rossi1/smart-pack/common/decorator"
	"github.com/rossi1/smart-pack/domain"
)

type SetPackSizesCommand struct {
	Sizes []domain.SmartPack
}

//go:generate mockgen -package=command -destination=set_pack_sizes.mock.go -source=set_pack_sizes.go
type SetPackSizesRepository interface {
	SetPackSizes(ctx context.Context, sizes []domain.SmartPack) error
}

type SetPackSizesHandler decorator.CommandHandler[*SetPackSizesCommand]

type setPackSizesHandler struct {
	repo SetPackSizesRepository
}

func NewSetPackSizesHandler(repo SetPackSizesRepository) SetPackSizesHandler {
	return decorator.ApplyCommandDecorators[*SetPackSizesCommand](&setPackSizesHandler{
		repo: repo,
	})
}

func (h *setPackSizesHandler) Handle(ctx context.Context, cmd *SetPackSizesCommand) error {
	return h.repo.SetPackSizes(ctx, cmd.Sizes)
}

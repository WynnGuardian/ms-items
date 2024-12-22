package usecase

import (
	"context"

	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	"github.com/wynnguardian/ms-items/internal/infra/repository"
)

type GetRankCaseInput struct {
	ItemName string `json:"item_name"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

type GetRankCase struct {
	Uow uow.UowInterface
}

func NewGetRankCase(uow uow.UowInterface) *GetRankCase {
	return &GetRankCase{
		Uow: uow,
	}
}

func (u *GetRankCase) Execute(ctx context.Context, in GetRankCaseInput) response.WGResponse {

	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {
		itemRepo := repository.GetAuthenticatedItemRepository(ctx, uow)

		rank, err := itemRepo.GetRank(ctx, in.ItemName, in.Page, in.Limit)
		if err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.New(200, "", rank)
	})
}

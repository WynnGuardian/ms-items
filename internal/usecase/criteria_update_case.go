package usecase

import (
	"context"

	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	"github.com/wynnguardian/ms-items/internal/infra/decoder/parser"
	"github.com/wynnguardian/ms-items/internal/infra/repository"
)

type CriteriaUpdateCaseInput struct {
	ItemName   string `json:"item_name"`
	CriteriaId string `json:"criteria_id"`
	Value      int    `json:"value"`
}

type CriteriaUpdateCase struct {
	Uow uow.UowInterface
}

func NewCriteriaUpdateCase(uow uow.UowInterface) *CriteriaUpdateCase {
	return &CriteriaUpdateCase{
		Uow: uow,
	}
}

func (u *CriteriaUpdateCase) Execute(ctx context.Context, in CriteriaUpdateCaseInput) response.WGResponse {

	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {
		criteriaRepo := repository.GetItemCriteriaRepository(ctx, uow)

		if _, ok := parser.NameToId[in.CriteriaId]; !ok {
			return response.ErrCriteriaNotFound
		}

		err := criteriaRepo.UpdateOne(ctx, in.ItemName, in.CriteriaId, float64(in.Value/100))
		if err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.Ok
	})
}

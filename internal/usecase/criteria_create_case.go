package usecase

import (
	"context"

	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	"github.com/wynnguardian/ms-items/internal/infra/decoder/parser"
	"github.com/wynnguardian/ms-items/internal/infra/repository"
)

type CreateCriteriaCaseInput struct {
	ItemName   string  `json:"item_name"`
	CriteriaId string  `json:"criteria_id"`
	Default    float64 `json:"default"`
}

type CreateCriteriaCase struct {
	Uow uow.UowInterface
}

func NewCreateCriteriaCase(uow uow.UowInterface) *CreateCriteriaCase {
	return &CreateCriteriaCase{
		Uow: uow,
	}
}

func (u *CreateCriteriaCase) Execute(ctx context.Context, in CreateCriteriaCaseInput) response.WGResponse {

	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {
		criteriaRepo := repository.GetItemCriteriaRepository(ctx, uow)

		if _, ok := parser.NameToId[in.CriteriaId]; !ok {
			return response.ErrCriteriaNotFound
		}

		err := criteriaRepo.Create(ctx, in.ItemName, in.CriteriaId, in.Default)
		if err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.Ok
	})
}

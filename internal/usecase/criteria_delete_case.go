package usecase

import (
	"context"
	"fmt"

	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	"github.com/wynnguardian/ms-items/internal/infra/decoder/parser"
	"github.com/wynnguardian/ms-items/internal/infra/repository"
)

type DeleteCriteriaCaseInput struct {
	ItemName   string  `json:"item_name"`
	CriteriaId string  `json:"criteria_id"`
	Default    float64 `json:"default"`
}

type DeleteCriteriaCase struct {
	Uow uow.UowInterface
}

func NewDeleteCriteriaCase(uow uow.UowInterface) *DeleteCriteriaCase {
	return &DeleteCriteriaCase{
		Uow: uow,
	}
}

func (u *DeleteCriteriaCase) Execute(ctx context.Context, in DeleteCriteriaCaseInput) response.WGResponse {

	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {
		criteriaRepo := repository.GetItemCriteriaRepository(ctx, uow)

		if _, ok := parser.NameToId[in.CriteriaId]; !ok {
			return response.ErrCriteriaNotFound
		}
		fmt.Println("OPA")
		err := criteriaRepo.Delete(ctx, in.ItemName, in.CriteriaId)
		if err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.Ok
	})
}

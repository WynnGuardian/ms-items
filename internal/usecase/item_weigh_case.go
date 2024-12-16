package usecase

import (
	"context"
	"net/http"

	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	"github.com/wynnguardian/ms-items/internal/infra/decoder"
	"github.com/wynnguardian/ms-items/internal/infra/decoder/parser"
	"github.com/wynnguardian/ms-items/internal/infra/repository"
	"github.com/wynnguardian/ms-items/internal/util"
)

type ItemWeighCaseInput struct {
	ItemUTF16 string `json:"item_utf16"`
}

type ItemWeighCaseOutput struct {
	Weight   float64              `json:"weight"`
	Item     *entity.ItemInstance `json:"item"`
	Criteria *entity.ItemCriteria `json:"criteria"`
}

type ItemWeighCase struct {
	Uow uow.UowInterface
}

func NewItemWeighCase(uow uow.UowInterface) *ItemWeighCase {
	return &ItemWeighCase{
		Uow: uow,
	}
}

func (u *ItemWeighCase) Execute(ctx context.Context, in ItemWeighCaseInput) response.WGResponse {
	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {

		wItemRepo := repository.GetWynnItemRepository(ctx, uow)
		criteriaRepo := repository.GetItemCriteriaRepository(ctx, uow)

		decoded, err := decoder.NewItemDecoder(in.ItemUTF16).Decode()
		if err != nil {
			return response.ErrInvalidItem
		}

		wynnItem, err := wItemRepo.Find(ctx, decoded.Name)
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrWynnItemNotFound)
		}

		criteria, err := criteriaRepo.Find(ctx, decoded.Name)
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrCriteriaNotFound)
		}

		parsed, err := parser.ParseDecodedItem(ctx, decoded, wynnItem)
		if err != nil {
			return response.ErrInternalServerErr(err)
		}

		weight := util.WeightItem(parsed, criteria)

		return response.New(http.StatusOK, "", ItemWeighCaseOutput{
			Weight:   weight,
			Item:     parsed,
			Criteria: criteria,
		})
	})
}

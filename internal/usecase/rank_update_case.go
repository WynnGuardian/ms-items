package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	"github.com/wynnguardian/common/utils"
	"github.com/wynnguardian/ms-items/internal/infra/decoder/parser"
	"github.com/wynnguardian/ms-items/internal/infra/repository"
	"github.com/wynnguardian/ms-items/internal/util"
)

type RankUpdateCaseInput struct {
	ItemName string `json:"item_name"`
}

type RankUpdateCase struct {
	Uow uow.UowInterface
}

func NewRankUpdateCase(uow uow.UowInterface) *RankUpdateCase {
	return &RankUpdateCase{
		Uow: uow,
	}
}

func (u *RankUpdateCase) Execute(ctx context.Context, in RankUpdateCaseInput) response.WGResponse {

	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {

		wItemRepo := repository.GetWynnItemRepository(ctx, uow)
		authenticatedRepo := repository.GetAuthenticatedItemRepository(ctx, uow)
		criteriaRepo := repository.GetItemCriteriaRepository(ctx, uow)

		wItem, err := wItemRepo.Find(ctx, in.ItemName)
		if err != nil {
			return utils.NotFoundOrInternalErr(err, response.ErrWynnItemNotFound)
		}

		criteria, err := criteriaRepo.Find(ctx, in.ItemName)
		if err != nil {
			return utils.NotFoundOrInternalErr(err, response.ErrCriteriaNotFound)
		}

		items, err := authenticatedRepo.FindAllWithItem(ctx, in.ItemName)
		if err != nil && err == sql.ErrNoRows {
			return response.Ok
		}

		if err != nil {
			return response.ErrInternalServerErr(err)
		}

		for _, item := range items {
			itemInstance, err := parser.ParseAuthenticatedItem(ctx, wItem, item)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			weight := util.WeightItem(itemInstance, criteria)
			item.LastRanked = time.Now()
			item.Weight = weight * 100
			if err := authenticatedRepo.Update(ctx, item); err != nil {
				fmt.Printf("Error while update item with tracking code %s\n: %s", item.TrackingCode, err.Error())
				return response.ErrInternalServerErr(err)
			}
		}

		items, err = authenticatedRepo.GetRank(ctx, in.ItemName, 1, 5000)
		if err != nil {
			return response.ErrInternalServerErr(err)
		}

		for i, item := range items {
			item.Position = i + 1
			if err := authenticatedRepo.Update(ctx, item); err != nil {
				fmt.Printf("Error while update item position with tracking code %s\n: %s", item.TrackingCode, err.Error())
				return response.ErrInternalServerErr(err)
			}
		}

		return response.Ok

	})

}

package usecase

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	"github.com/wynnguardian/common/utils"
	"github.com/wynnguardian/ms-items/internal/infra/decoder"
	"github.com/wynnguardian/ms-items/internal/infra/decoder/parser"
	"github.com/wynnguardian/ms-items/internal/infra/repository"
	"github.com/wynnguardian/ms-items/internal/util"
)

type AuthenticateItemCaseInput struct {
	ItemUTF16  string `json:"item_utf16"`
	MCOwnerUID string `json:"owner_mc_uid"`
	DCOwnerUID string `json:"owner_dc_uid"`
	Public     bool   `json:"public_info"`
	Force      bool   `json:"force"`
}

type AuthenticateItemCaseOutput struct {
	TrackingCode string                    `json:"tracking_code"`
	WynnItem     *entity.WynnItem          `json:"wynn_item"`
	Weight       float64                   `json:"weight"`
	Item         *entity.AuthenticatedItem `json:"item"`
}

type AuthenticateItemCase struct {
	Uow uow.UowInterface
}

func NewAuthenticatetemCase(uow uow.UowInterface) *AuthenticateItemCase {
	return &AuthenticateItemCase{
		Uow: uow,
	}
}

func (u *AuthenticateItemCase) Execute(ctx context.Context, in AuthenticateItemCaseInput) response.WGResponse {

	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {

		wItemRepo := repository.GetWynnItemRepository(ctx, uow)
		criteriaRepo := repository.GetItemCriteriaRepository(ctx, uow)
		authRepo := repository.GetAuthenticatedItemRepository(ctx, uow)

		tCode := utils.GenAuthId()
		id := util.GenNanoId(24)

		d := decoder.NewItemDecoder(in.ItemUTF16)

		var bytes string
		for _, num := range d.Reader.Data {
			bytes += fmt.Sprintf("%d", num)
		}

		_, err := authRepo.FindWithBytes(ctx, bytes)
		if err == nil && !in.Force {
			return response.New[any](http.StatusBadRequest, "Item with same bytes found. May it be duplicated? Use FORCE = TRUE to force authentication.", nil)
		}

		decoded, err := d.Decode()

		if err != nil {
			return response.ErrInvalidItem
		}

		expected, err := wItemRepo.Find(ctx, decoded.Name)
		if err != nil {
			return response.ErrWynnItemNotFound
		}

		criteria, err := criteriaRepo.Find(ctx, decoded.Name)
		if err != nil {
			return response.ErrCriteriaNotFound
		}

		item, err := parser.ParseDecodedItem(ctx, decoded, expected)
		if err != nil {
			return response.ErrInternalServerErr(err)
		}

		if ok := util.HasAllCriterias(item, criteria); !ok {
			return response.New[any](http.StatusBadRequest, "Item does not have all mandatory criteria", nil)
		}

		weight := util.WeightItem(item, criteria)

		i := &entity.AuthenticatedItem{
			Id:           id,
			Bytes:        bytes,
			Position:     9999,
			Item:         expected.Name,
			OwnerMC:      in.MCOwnerUID,
			OwnerDC:      in.DCOwnerUID,
			Stats:        item.Stats,
			Weight:       weight * 100,
			LastRanked:   time.Now(),
			PublicOwner:  in.Public,
			TrackingCode: tCode,
		}

		if err = authRepo.Create(ctx, i); err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.New(http.StatusOK, "", AuthenticateItemCaseOutput{
			TrackingCode: tCode,
			Item:         i,
			Weight:       weight,
			WynnItem:     expected,
		})

	})

}

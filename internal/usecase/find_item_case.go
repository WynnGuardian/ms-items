package usecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	"github.com/wynnguardian/ms-items/internal/infra/repository"
	"github.com/wynnguardian/ms-items/internal/util"
)

type FindItemCaseInput struct {
	TrackingCode string `json:"tracking_code"`
}

type FindItemCase struct {
	Uow uow.UowInterface
}

func NewFindItemCase(uow uow.UowInterface) *FindItemCase {
	return &FindItemCase{
		Uow: uow,
	}
}

func (u *FindItemCase) Execute(ctx context.Context, in FindItemCaseInput) response.WGResponse {

	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {
		itemRepo := repository.GetAuthenticatedItemRepository(ctx, uow)

		cr, err := itemRepo.Find(ctx, in.TrackingCode)
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrSurveyNotFound)
		}

		fmt.Println(cr)

		return response.New(http.StatusOK, "", cr)
	})

}

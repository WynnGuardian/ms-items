package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	"github.com/wynnguardian/ms-items/internal/usecase"
)

func WeightItem(ctx *gin.Context) response.WGResponse {
	input := usecase.ItemWeighCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewItemWeighCase(uow.Current()).Execute(ctx, input)
}

func AuthItem(ctx *gin.Context) response.WGResponse {
	input := usecase.AuthenticateItemCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewAuthenticatetemCase(uow.Current()).Execute(ctx, input)
}

func FindCriteria(ctx *gin.Context) response.WGResponse {
	input := usecase.FindCriteriaCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewFindCriteriaCase(uow.Current()).Execute(ctx, input)
}

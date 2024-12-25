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

func UpdateRank(ctx *gin.Context) response.WGResponse {
	input := usecase.RankUpdateCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewRankUpdateCase(uow.Current()).Execute(ctx, input)
}

func GetRank(ctx *gin.Context) response.WGResponse {
	input := usecase.GetRankCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewGetRankCase(uow.Current()).Execute(ctx, input)
}

func CreateCriteria(ctx *gin.Context) response.WGResponse {
	input := usecase.CreateCriteriaCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewCreateCriteriaCase(uow.Current()).Execute(ctx, input)
}

func DeleteCriteria(ctx *gin.Context) response.WGResponse {
	input := usecase.DeleteCriteriaCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewDeleteCriteriaCase(uow.Current()).Execute(ctx, input)
}

func UpdateCriteria(ctx *gin.Context) response.WGResponse {
	input := usecase.CriteriaUpdateCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewCriteriaUpdateCase(uow.Current()).Execute(ctx, input)
}

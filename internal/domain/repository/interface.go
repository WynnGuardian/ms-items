package repository

import (
	"context"

	"github.com/wynnguardian/common/entity"
)

type RepositoryInterface interface {
}

type AuthenticatedItemRepositoryInterface interface {
	FindFirst(ctx context.Context, id string) (*entity.AuthenticatedItem, error)
	FindAllWithItem(ctx context.Context, name string) ([]*entity.AuthenticatedItem, error)
	Update(ctx context.Context, item *entity.AuthenticatedItem) error
	Create(ctx context.Context, item *entity.AuthenticatedItem) error
	GetRank(ctx context.Context, itemName string, page, limit int) ([]*entity.AuthenticatedItem, error)
}

type WynnItemRepositoryInterface interface {
	Find(ctx context.Context, name string) (*entity.WynnItem, error)
}

type CriteriaRepositoryInterface interface {
	Find(ctx context.Context, name string) (*entity.ItemCriteria, error)
	Update(ctx context.Context, crit *entity.ItemCriteria) error
}

type GenRepositoryInterface interface {
	GenItemDB(ctx context.Context)
	GenDefaultScales(ctx context.Context) error
}

package repository

import (
	"context"
	"errors"

	"github.com/victorbetoni/wynnguardian/ms-items/internal/domain/repository"
	"github.com/victorbetoni/wynnguardian/ms-items/internal/infra/db"
	"github.com/wynnguardian/common/uow"
)

var ErrQueriesNotSet = errors.New("queries not set")

type Repository struct {
	*db.Queries
}

func (r *Repository) SetQuery(q *db.Queries) {
	r.Queries = q
}

func (r *Repository) Validate() error {
	if r.Queries == nil {
		return ErrQueriesNotSet
	}
	return nil
}

func GetWynnItemRepository(ctx context.Context, u *uow.Uow) repository.WynnItemRepositoryInterface {
	return getRepository[repository.WynnItemRepositoryInterface](ctx, u, "WynnItemRepository")
}

func GetAuthenticatedItemRepository(ctx context.Context, u *uow.Uow) repository.AuthenticatedItemRepositoryInterface {
	return getRepository[repository.AuthenticatedItemRepositoryInterface](ctx, u, "AuthenticatedItemRepository")
}

func GetGenRepository(ctx context.Context, u *uow.Uow) repository.GenRepositoryInterface {
	return getRepository[repository.GenRepositoryInterface](ctx, u, "GenRepository")
}

func GetItemCriteriaRepository(ctx context.Context, u *uow.Uow) repository.CriteriaRepositoryInterface {
	return getRepository[repository.CriteriaRepositoryInterface](ctx, u, "CriteriaRepository")
}

func getRepository[T repository.RepositoryInterface](ctx context.Context, u *uow.Uow, name string) T {
	repo, err := u.GetRepository(ctx, name)
	if err != nil {
		panic(err)
	}
	return repo.(T)
}

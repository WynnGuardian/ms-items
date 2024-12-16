package repository

import (
	"context"
	"database/sql"

	"github.com/victorbetoni/go-streams/streams"
	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/ms-items/internal/infra/db"
)

type AuthenticatedItemRepository struct {
	dbConn *sql.DB
	*db.Queries
}

func NewAuthenticatedItemRepository(dbConn *sql.DB) *AuthenticatedItemRepository {
	return &AuthenticatedItemRepository{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}

func (r *AuthenticatedItemRepository) Find(ctx context.Context, id string) (*entity.AuthenticatedItem, error) {
	i, err := r.Queries.FindAuthenticatedItem(ctx, db.FindAuthenticatedItemParams{
		Code: id,
	})

	if err != nil {
		return nil, err
	}

	stats, err := r.Queries.FindAuthenticatedItemStats(ctx, db.FindAuthenticatedItemStatsParams{
		Code: id,
	})

	if err != nil {
		return nil, err
	}

	st := make(map[string]int, 0)
	for _, s := range stats {
		st[s.Statid] = int(s.Value)
	}

	return &entity.AuthenticatedItem{
		Id:           i.ID,
		Item:         i.Itemname,
		OwnerMC:      i.Ownermcuuid,
		OwnerDC:      i.Owneruserid,
		Stats:        st,
		Position:     int(i.Position),
		LastRanked:   i.Lastranked,
		PublicOwner:  int(i.Ownerpublic) != 0,
		TrackingCode: i.Trackingcode,
		Bytes:        i.Bytes,
	}, nil
}

func (r *AuthenticatedItemRepository) FindAllWithItem(ctx context.Context, name string) ([]*entity.AuthenticatedItem, error) {

	items, err := r.Queries.FindWynnItemAuthenticatedItems(ctx, name)
	if err != nil {
		return nil, err
	}

	return *streams.Map(streams.StreamOf(items...), func(t db.WgAuthenticateditem) *entity.AuthenticatedItem {

		statsMap := make(map[string]int, 0)
		if st, err := r.Queries.FindAuthenticatedItemStats(ctx, db.FindAuthenticatedItemStatsParams{
			Code: t.ID,
		}); err == nil {
			for _, s := range st {
				statsMap[s.Statid] = int(s.Value)
			}
		}

		return &entity.AuthenticatedItem{
			Id:           t.ID,
			Item:         t.Itemname,
			OwnerMC:      t.Ownermcuuid,
			OwnerDC:      t.Owneruserid,
			Stats:        statsMap,
			Position:     int(t.Position),
			LastRanked:   t.Lastranked,
			PublicOwner:  t.Ownerpublic != 0,
			TrackingCode: t.Trackingcode,
			Bytes:        t.Bytes,
		}
	}).ToSlice(), nil
}

func (r *AuthenticatedItemRepository) Create(ctx context.Context, item *entity.AuthenticatedItem) error {

	p := 1
	if !item.PublicOwner {
		p = 0
	}

	return r.Queries.CreateAuthenticatedItem(ctx, db.CreateAuthenticatedItemParams{
		ID:           item.Id,
		Lastranked:   item.LastRanked,
		Itemname:     item.Item,
		Ownermcuuid:  item.OwnerMC,
		Owneruserid:  item.OwnerDC,
		Trackingcode: item.TrackingCode,
		Ownerpublic:  int32(p),
		Position:     int32(item.Position),
		Bytes:        item.Bytes,
	})
}

func (r *AuthenticatedItemRepository) Update(ctx context.Context, item *entity.AuthenticatedItem) error {

	p := 1
	if !item.PublicOwner {
		p = 0
	}

	return r.Queries.UpdateAuthenticatedItem(ctx, db.UpdateAuthenticatedItemParams{
		ID:           item.Id,
		Lastranked:   item.LastRanked,
		Ownermcuuid:  item.OwnerMC,
		Owneruserid:  item.OwnerDC,
		Trackingcode: item.TrackingCode,
		Ownerpublic:  int32(p),
		Position:     int32(item.Position),
		Bytes:        item.Bytes,
	})
}

package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/victorbetoni/go-streams/streams"
	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/common/uow"
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
		Code: i.ID,
	})

	if err != nil {
		return nil, err
	}

	st := make(map[string]int, 0)
	for _, s := range stats {
		st[s.Statid] = int(s.Value)
	}

	w, err := GetWynnItemRepository(ctx, uow.Current()).Find(ctx, i.Itemname)
	if err != nil {
		return nil, err
	}

	return &entity.AuthenticatedItem{
		Id:           i.ID,
		WynnItem:     w,
		Position:     int(i.Position),
		Item:         i.Itemname,
		OwnerMC:      i.Ownermcuuid,
		OwnerDC:      i.Owneruserid,
		Stats:        st,
		Weight:       i.Weight,
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
			Weight:       t.Weight,
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

	err := r.Queries.CreateAuthenticatedItem(ctx, db.CreateAuthenticatedItemParams{
		ID:           item.Id,
		Lastranked:   item.LastRanked,
		Itemname:     item.Item,
		Ownermcuuid:  item.OwnerMC,
		Owneruserid:  item.OwnerDC,
		Position:     int32(item.Position),
		Trackingcode: item.TrackingCode,
		Ownerpublic:  int32(p),
		Weight:       item.Weight,
		Bytes:        item.Bytes,
	})

	if err != nil {
		return err
	}

	for st, v := range item.Stats {
		if err := r.Queries.CreateAuthenticatedItemStat(ctx, db.CreateAuthenticatedItemStatParams{
			Itemid: item.Id,
			Statid: st,
			Value:  int32(v),
		}); err != nil {
			return err
		}
	}

	return nil
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
		Weight:       item.Weight,
		Bytes:        item.Bytes,
		Position:     int32(item.Position),
	})
}

func (r *AuthenticatedItemRepository) GetRank(ctx context.Context, itemName string, page, limit int) ([]*entity.AuthenticatedItem, error) {
	i, err := r.Queries.RankAuthenticatedItems(ctx, db.RankAuthenticatedItemsParams{
		Itemname: itemName,
		Limit:    int32(limit),
		Offset:   int32(page-1) * int32(limit),
	})

	wRepo := GetWynnItemRepository(ctx, uow.Current())

	if err != nil {
		return nil, err
	}

	return *streams.Map(streams.StreamOf(i...), func(i db.WgAuthenticateditem) *entity.AuthenticatedItem {

		statsMap := make(map[string]int, 0)
		if st, err := r.Queries.FindAuthenticatedItemStats(ctx, db.FindAuthenticatedItemStatsParams{
			Code: i.ID,
		}); err == nil {
			for _, s := range st {
				statsMap[s.Statid] = int(s.Value)
			}
		}

		it := &entity.AuthenticatedItem{
			Id:           i.ID,
			Item:         i.Itemname,
			OwnerMC:      i.Ownermcuuid,
			OwnerDC:      i.Owneruserid,
			Stats:        statsMap,
			Position:     int(i.Position),
			Weight:       i.Weight,
			LastRanked:   i.Lastranked,
			PublicOwner:  int(i.Ownerpublic) != 0,
			TrackingCode: i.Trackingcode,
			Bytes:        i.Bytes,
		}

		wItem, err := wRepo.Find(ctx, i.Itemname)
		if err != nil {
			fmt.Println(err.Error())
			return it
		}
		it.WynnItem = wItem
		return it
	}).ToSlice(), nil
}

func (r *AuthenticatedItemRepository) FindWithBytes(ctx context.Context, bytes string) (*entity.AuthenticatedItem, error) {
	v, err := r.Queries.FindWithBytes(ctx, bytes)
	if err != nil {
		return nil, err
	}

	w, err := GetWynnItemRepository(ctx, uow.Current()).Find(ctx, v.Itemname)
	if err != nil {
		return nil, err
	}
	return &entity.AuthenticatedItem{
		Id:           v.ID,
		WynnItem:     w,
		Item:         v.Itemname,
		Weight:       v.Weight,
		Position:     int(v.Position),
		OwnerMC:      v.Ownermcuuid,
		OwnerDC:      v.Owneruserid,
		Stats:        nil,
		LastRanked:   v.Lastranked,
		PublicOwner:  v.Ownerpublic == 1,
		TrackingCode: v.Trackingcode,
		Bytes:        v.Bytes,
	}, nil
}

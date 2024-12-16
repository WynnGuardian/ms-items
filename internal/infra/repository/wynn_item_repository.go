package repository

import (
	"context"
	"database/sql"

	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/ms-items/internal/infra/db"
)

type WynnItemRepository struct {
	dbConn *sql.DB
	*db.Queries
}

func NewWynnItemRepository(dbConn *sql.DB) *WynnItemRepository {
	return &WynnItemRepository{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}

func (r *WynnItemRepository) Find(ctx context.Context, name string) (*entity.WynnItem, error) {

	item, err := r.Queries.FindWynnItem(ctx, name)
	if err != nil {
		return nil, err
	}

	stats, err := r.Queries.FindWynnItemStats(ctx, name)
	if err != nil {
		return nil, err
	}

	generic := &entity.WynnItem{
		Name: item.Name,
		Requirements: entity.Requirements{
			CombatLevel:  int(item.Reqlevel),
			Strenght:     int(item.Reqstrenght),
			Dexterity:    int(item.Reqdexterity),
			Intelligence: int(item.Reqintelligence),
			Defence:      int(item.Reqdefence),
			Agility:      int(item.Reqagility),
		},
	}

	statsDef := make(map[string]entity.Stat, 0)
	for _, stat := range stats {
		statsDef[stat.Statid] = entity.Stat{
			Id:      stat.Statid,
			Minimum: int(stat.Lower),
			Maximum: int(stat.Upper),
		}
	}

	generic.Stats = statsDef
	return generic, nil
}

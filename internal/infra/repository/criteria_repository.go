package repository

import (
	"context"
	"database/sql"

	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/ms-items/internal/infra/db"
)

type CriteriaRepository struct {
	dbConn *sql.DB
	*db.Queries
}

func NewCriteriaRepository(dbConn *sql.DB) *CriteriaRepository {
	return &CriteriaRepository{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}

func (r *CriteriaRepository) Find(ctx context.Context, name string) (*entity.ItemCriteria, error) {
	mods, err := r.Queries.FindItemCriterias(ctx, name)

	if err != nil {
		return nil, err
	}

	modifiers := make(map[string]float64, 0)
	for _, m := range mods {
		modifiers[m.Statid] = float64(m.Value)
	}

	return &entity.ItemCriteria{
		Item:      name,
		Modifiers: modifiers,
	}, nil

}

func (c *CriteriaRepository) Update(ctx context.Context, crit *entity.ItemCriteria) error {
	for st, val := range crit.Modifiers {
		if err := c.Queries.UpdateCriteria(ctx, db.UpdateCriteriaParams{
			Value:    val,
			Itemname: crit.Item,
			Statid:   st,
		}); err != nil {
			return err
		}
	}
	return nil
}

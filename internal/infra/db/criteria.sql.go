// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: criteria.sql

package db

import (
	"context"
)

const clearCriteriaTable = `-- name: ClearCriteriaTable :exec
DELETE FROM WG_Criteria
`

func (q *Queries) ClearCriteriaTable(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, clearCriteriaTable)
	return err
}

const createCriteria = `-- name: CreateCriteria :exec
INSERT INTO WG_Criteria (ItemName, StatId, Value) VALUES (?,?,?)
`

type CreateCriteriaParams struct {
	Itemname string  `json:"itemname"`
	Statid   string  `json:"statid"`
	Value    float64 `json:"value"`
}

func (q *Queries) CreateCriteria(ctx context.Context, arg CreateCriteriaParams) error {
	_, err := q.db.ExecContext(ctx, createCriteria, arg.Itemname, arg.Statid, arg.Value)
	return err
}

const findItemCriterias = `-- name: FindItemCriterias :many
SELECT itemname, statid, value FROM WG_Criteria WHERE ItemName = ?
`

func (q *Queries) FindItemCriterias(ctx context.Context, itemname string) ([]WgCriterium, error) {
	rows, err := q.db.QueryContext(ctx, findItemCriterias, itemname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []WgCriterium
	for rows.Next() {
		var i WgCriterium
		if err := rows.Scan(&i.Itemname, &i.Statid, &i.Value); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCriteria = `-- name: UpdateCriteria :exec
UPDATE WG_Criteria SET Value = ? WHERE ItemName = ? AND StatId = ?
`

type UpdateCriteriaParams struct {
	Value    float64 `json:"value"`
	Itemname string  `json:"itemname"`
	Statid   string  `json:"statid"`
}

func (q *Queries) UpdateCriteria(ctx context.Context, arg UpdateCriteriaParams) error {
	_, err := q.db.ExecContext(ctx, updateCriteria, arg.Value, arg.Itemname, arg.Statid)
	return err
}
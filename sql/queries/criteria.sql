-- name: UpdateCriteria :exec
UPDATE WG_Criteria SET Value = ? WHERE ItemName = ? AND StatId = ?;

-- name: FindItemCriterias :many
SELECT * FROM WG_Criteria WHERE ItemName = ? ORDER BY Value DESC;

-- name: CreateCriteria :exec
INSERT INTO WG_Criteria (ItemName, StatId, Value) VALUES (?,?,?);

-- name: ClearCriteriaTable :exec
DELETE FROM WG_Criteria;

-- name: DeleteCriteria :exec
DELETE FROM WG_Criteria WHERE ItemName = ? AND StatId = ?;
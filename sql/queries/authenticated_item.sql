-- name: CreateAuthenticatedItem :exec
INSERT INTO WG_AuthenticatedItems (Id, LastRanked, ItemName, OwnerMCUUID, OwnerUserId, Weight, TrackingCode, OwnerPublic, Bytes) VALUES (?,?,?,?,?,?,?,?,?);

-- name: FindAuthenticatedItemStats :many
SELECT * FROM WG_AuthenticatedItemStats WHERE ItemId = sqlc.arg(code) OR ItemId = sqlc.arg(code);

-- name: FindAuthenticatedItem :one
SELECT * FROM WG_AuthenticatedItems WHERE Id = sqlc.arg(code) OR TrackingCode = sqlc.arg(code);

-- name: UpdateAuthenticatedItem :exec
UPDATE WG_AuthenticatedItems SET LastRanked = ?, OwnerMCUUID = ?, OwnerUserId = ?,  Weight = ?, OwnerPublic = ?, Bytes = ? WHERE Id = ? OR TrackingCode = ?;

-- name: FindAllAuthenticatedItemNames :many
SELECT DISTINCT ItemName FROM WG_AuthenticatedItems;

-- name: FindWynnItemAuthenticatedItems :many
SELECT * FROM WG_AuthenticatedItems WHERE ItemName = ?;

-- name: RankAuthenticatedItems :many
SELECT * FROM WG_AuthenticatedItems WHERE ItemName = ? ORDER BY Weight DESC LIMIT ? OFFSET ?;
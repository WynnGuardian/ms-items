-- name: FindWynnItem :one
SELECT * FROM WG_WynnItems WHERE Name = ?;

-- name: FindWynnItemStats :many
SELECT * FROM WG_WynnItemStats WHERE ItemName = ?;

-- name: CreateWynnItem :exec
INSERT INTO WG_WynnItems(Name, Sprite, ReqLevel, ReqStrenght, ReqAgility, ReqDefence, ReqIntelligence, ReqDexterity) VALUES (?,?,?,?,?,?,?,?);

-- name: CreateWynnItemStat :exec
INSERT INTO WG_WynnItemStats (ItemName, StatId, Lower, Upper) VALUES (?,?,?,?);

-- name: ClearWynnItemsTable :exec
DELETE FROM WG_WynnItems;

-- name: ClearWynnItemStats :exec
DELETE FROM WG_WynnItemStats;
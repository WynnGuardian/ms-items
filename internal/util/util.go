package util

import (
	"database/sql"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/wynnguardian/common/response"
)

func GenNanoId(size int) string {
	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", size)
	if err != nil {
		panic(err)
	}
	return id
}

func NotFoundOrInternalErr(err error, notFound response.WGResponse) response.WGResponse {
	if err == sql.ErrNoRows {
		return notFound
	}
	return response.ErrInternalServerErr(err)
}

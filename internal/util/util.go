package util

import gonanoid "github.com/matoous/go-nanoid/v2"

func GenNanoId(size int) string {
	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", size)
	if err != nil {
		panic(err)
	}
	return id
}

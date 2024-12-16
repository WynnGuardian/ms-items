package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/wynnguardian/common/uow"
	util "github.com/wynnguardian/common/utils"
	"github.com/wynnguardian/ms-items/internal/infra/db"
	"github.com/wynnguardian/ms-items/internal/infra/decoder/parser"
	"github.com/wynnguardian/ms-items/internal/infra/http/router"
	"github.com/wynnguardian/ms-items/internal/infra/repository"
)

func main() {

	ctx := context.Background()

	parser.LoadIdTable()
	util.Must(godotenv.Load(".env"))
	db := util.MustVal(sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(mysql:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PW"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))+"?parseTime=true&loc=America%2FSao_Paulo"))
	util.Must(db.Ping())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	uow := util.MustVal(uow.NewUow(ctx, db))

	registerRepositories(uow)

	/*uow.Do(ctx, func(uow *u.Uow) response.WGResponse {
		repo := repository.GetGenRepository(ctx, uow)
		repo.GenDefaultScales(ctx)
		repo.GenItemDB(ctx)
		return response.WGResponse{
			Status: 200,
		}
	})*/

	defer db.Close()

	r := router.Build()

	err := r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatalf("Couldn't start HTTP server: %s\n", err.Error())
		return
	}
	fmt.Println("Listening on port ", os.Getenv("PORT"))
}

func registerRepositories(uow *uow.Uow) {
	uow.Register("WynnItemRepository", func(tx *sql.Tx) interface{} {
		repo := repository.NewWynnItemRepository(uow.Db)
		repo.Queries = db.New(tx)
		return repo
	})

	uow.Register("AuthenticatedItemRepository", func(tx *sql.Tx) interface{} {
		repo := repository.NewAuthenticatedItemRepository(uow.Db)
		repo.Queries = db.New(tx)
		return repo
	})

	uow.Register("CriteriaRepository", func(tx *sql.Tx) interface{} {
		repo := repository.NewCriteriaRepository(uow.Db)
		repo.Queries = db.New(tx)
		return repo
	})

	uow.Register("GenRepository", func(tx *sql.Tx) interface{} {
		repo := repository.NewGenRepository(uow.Db)
		repo.Queries = db.New(tx)
		return repo
	})

}

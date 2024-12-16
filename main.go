package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/wynnguardian/common/uow"
	"github.com/wynnguardian/common/utils"
	"github.com/wynnguardian/ms-items/internal/domain/config"
	"github.com/wynnguardian/ms-items/internal/infra/db"
	"github.com/wynnguardian/ms-items/internal/infra/decoder/parser"
	"github.com/wynnguardian/ms-items/internal/infra/http/router"
	"github.com/wynnguardian/ms-items/internal/infra/repository"
)

func main() {

	ctx := context.Background()

	parser.LoadIdTable()

	config.Load()

	db := utils.MustVal(sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.MainConfig.Private.DB.Username, config.MainConfig.Private.DB.Password, config.MainConfig.Private.DB.Hostname, config.MainConfig.Private.DB.Port, config.MainConfig.Private.DB.Database)+"?parseTime=true&loc=America%2FSao_Paulo"))
	utils.Must(db.Ping())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	uow := utils.MustVal(uow.NewUow(ctx, db))

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

	err := r.Run(fmt.Sprintf(":%d", config.MainConfig.Server.Port))
	if err != nil {
		log.Fatalf("Couldn't start HTTP server: %s\n", err.Error())
		return
	}
	fmt.Println("Listening on port ", config.MainConfig.Server.Port))
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

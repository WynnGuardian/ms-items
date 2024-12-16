package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/wynnguardian/ms-items/internal/infra/db"
)

type IncomingItem struct {
	Rarity            string `json:"rarity"`
	InternalName      string `json:"internalName"`
	LevelReq          int    `json:"levelReq"`
	StrReq            int    `json:"strReq"`
	AgiReq            int    `json:"agiReq"`
	IntReq            int    `json:"intReq"`
	DefReq            int    `json:"defReq"`
	DexReq            int    `json:"dexReq"`
	BaseHP            int    `json:"baseHP"`
	BaseEarthDef      int    `json:"baseEarthDef"`
	BaseAirDef        int    `json:"baseAirDef"`
	BaseThunderDef    int    `json:"baseThunderDef"`
	BaseWaterDef      int    `json:"baseWaterDef"`
	BaseFireDef       int    `json:"baseFireDef"`
	BaseDamage        []int  `json:"baseDamage"`
	BaseAirDamage     []int  `json:"baseAirDam"`
	BaseEarthDamage   []int  `json:"baseEarthDam"`
	BaseFireDamage    []int  `json:"baseFireDam"`
	BaseThunderDamage []int  `json:"baseThunderDam"`
	BaseWaterDamage   []int  `json:"baseWaterDam"`
	Identification    []struct {
		Stat string `json:"stat"`
		Max  int    `json:"max"`
		Min  int    `json:"min"`
	} `json:"identifications"`
}

type GenRepository struct {
	dbConn *sql.DB
	*db.Queries
}

func NewGenRepository(dbConn *sql.DB) *GenRepository {
	return &GenRepository{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}

func (g *GenRepository) GenDefaultScales(ctx context.Context) error {

	type Weights struct {
		Item   string `json:"item"`
		Weight []struct {
			Id     string  `json:"id"`
			Weight float32 `json:"weight"`
		} `json:"weights"`
	}

	file, err := os.Open("gen/default_weight.json")
	if err != nil {
		fmt.Printf("Error while opening scale json file: %s\n", err.Error())
		return err
	}
	defer file.Close()

	var weights []Weights

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&weights); err != nil {
		fmt.Printf("Error while decoding item json file: %s\n", err.Error())
		return err
	}
	fmt.Printf("Loaded %d scales\n", len(weights))

	for _, w := range weights {
		sum := float32(0)
		for _, m := range w.Weight {
			sum += m.Weight
		}
	}

	fmt.Println("Dropping previous criteria table...")
	err = g.Queries.ClearCriteriaTable(ctx)
	if err != nil {
		fmt.Printf("Error while dropping previous scale table: %s\n", err.Error())
		return err
	}

	for _, w := range weights {
		for _, m := range w.Weight {
			err := g.Queries.CreateCriteria(ctx, db.CreateCriteriaParams{
				Itemname: w.Item,
				Statid:   m.Id,
				Value:    float64(m.Weight),
			})
			if err != nil {
				fmt.Printf("error while inserting item scale for item %s and id %s: %s\n", w.Item, m.Id, err.Error())
				continue
			}
		}
	}
	return nil
}

func (g *GenRepository) GenItemDB(ctx context.Context) {
	fmt.Println("Dropping previous item stat table...")
	err := g.Queries.ClearWynnItemStats(ctx)
	if err != nil {
		fmt.Printf("Error while dropping previous item table: %s\n", err.Error())
		return
	}
	fmt.Println("Dropping previous item table...")
	err = g.Queries.ClearWynnItemsTable(ctx)
	if err != nil {
		fmt.Printf("Error while dropping previous item table: %s\n", err.Error())
		return
	}

	fmt.Println("Generating updated item json file...")
	cmd := exec.Command("node", "gen/json_sanitizer.js")
	_, err = cmd.Output()
	if err != nil {
		fmt.Printf("Error while generating item json: %s\n", err.Error())
		return
	}

	file, err := os.Open("gen/sanitized.json")
	if err != nil {
		fmt.Printf("Error while opening item json file: %s\n", err.Error())
		return
	}
	defer file.Close()

	var items []IncomingItem

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&items); err != nil {
		fmt.Printf("Error while decoding item json file: %s\n", err.Error())
		return
	}
	fmt.Printf("Loaded %d items\n", len(items))

	for _, item := range items {

		err := g.Queries.CreateWynnItem(ctx, db.CreateWynnItemParams{
			Name:            item.InternalName,
			Sprite:          "",
			Reqlevel:        int32(item.LevelReq),
			Reqstrenght:     int32(item.StrReq),
			Reqagility:      int32(item.AgiReq),
			Reqdefence:      int32(item.DefReq),
			Reqintelligence: int32(item.IntReq),
			Reqdexterity:    int32(item.DexReq),
		})
		if err != nil {
			fmt.Printf("Error while inserting generic item %s: %s\n", item.InternalName, err.Error())
			continue
		}
		for _, id := range item.Identification {
			err = g.Queries.CreateWynnItemStat(ctx, db.CreateWynnItemStatParams{
				Itemname: item.InternalName,
				Lower:    int32(id.Min),
				Upper:    int32(id.Max),
				Statid:   id.Stat,
			})
			if err != nil {
				fmt.Printf("Error while inserting generic item %s stat %s: %s\n", item.InternalName, id.Stat, err.Error())
				continue
			}
		}
	}
}

package parser

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var IdTable map[int]string

func LoadIdTable() {
	IdTable = make(map[int]string, 0)
	file, err := os.Open("id_table.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var d map[string]int
	if err := json.Unmarshal(data, &d); err != nil {
		log.Fatal(err)
	}

	for name, id := range d {
		IdTable[id] = name
	}

}

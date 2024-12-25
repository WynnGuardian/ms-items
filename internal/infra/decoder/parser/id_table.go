package parser

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var IdToName map[int]string
var NameToId map[string]int

func LoadIdTable() {
	IdToName = make(map[int]string, 0)
	NameToId = make(map[string]int, 0)
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

	NameToId = d
	for name, id := range d {
		IdToName[id] = name
	}

}

package entity

type Requirements struct {
	CombatLevel  int `json:"combat_level"`
	Strenght     int `json:"strenght"`
	Dexterity    int `json:"dexterity"`
	Intelligence int `json:"intelligence"`
	Defence      int `json:"defence"`
	Agility      int `json:"agility"`
}

type WynnItem struct {
	Name         string          `json:"name"`
	Requirements Requirements    `json:"requirements"`
	Stats        map[string]Stat `json:"stat"`
}

type Stat struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Minimum int    `json:"minimum"`
	Maximum int    `json:"maximum"`
}

type ItemInstance struct {
	Item     string         `json:"item_name"`
	Stats    map[string]int `json:"stats"`
	WynnItem *WynnItem      `json:"wynn_item"`
}

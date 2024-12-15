package entity

type ItemCriteria struct {
	Item      string             `json:"item"`
	Modifiers map[string]float64 `json:"modifiers"`
}

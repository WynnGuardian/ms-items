package util

import (
	"context"
	"strings"

	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/ms-items/internal/infra/decoder"
	"github.com/wynnguardian/ms-items/internal/infra/decoder/parser"
)

func WeightDecodedItem(ctx context.Context, decoded *decoder.DecodedItem, expected *entity.WynnItem, weightData *entity.ItemCriteria) (float64, error) {

	static, err := parser.ParseDecodedItem(ctx, decoded, expected)
	if err != nil {
		return 0, err
	}

	return WeightItem(static, weightData), nil

}

func WeightItem(item *entity.ItemInstance, criteria *entity.ItemCriteria) float64 {
	weight := 0.0

	for attribute, val := range item.Stats {
		norm := normalize(val, item.WynnItem.Stats[attribute].Minimum, item.WynnItem.Stats[attribute].Maximum, strings.Contains(attribute, "cost"))
		weight += norm * criteria.Modifiers[attribute]
	}

	return weight
}

func normalize(val, min, max int, reverse bool) float64 {
	diff := max - min
	offset := val - min
	if reverse {
		offset = max - val
	}
	return float64(offset) / float64(diff)
}

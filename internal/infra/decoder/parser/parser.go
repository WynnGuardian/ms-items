package parser

import (
	"context"
	"fmt"

	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/ms-items/internal/infra/decoder"
)

func ParseDecodedItem(ctx context.Context, decoded *decoder.DecodedItem, expected *entity.WynnItem) (*entity.ItemInstance, error) {

	canonicalIds := make(map[string]int, 0)
	for id, val := range decoded.Identifications {
		name, ok := IdTable[int(id)]
		if !ok {
			return nil, fmt.Errorf("unknow identification numeric id: \"%d\"", id)
		}
		canonicalIds[name] = val
	}

	for id := range canonicalIds {
		if _, ok := expected.Stats[id]; !ok {
			return nil, fmt.Errorf("ID %s does not exist in item \"%s\"", id, expected.Name)
		}
	}

	return &entity.ItemInstance{
		Item:     expected.Name,
		Stats:    canonicalIds,
		WynnItem: expected,
	}, nil

}

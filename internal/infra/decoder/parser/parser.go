package parser

import (
	"context"
	"fmt"

	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/ms-items/internal/infra/decoder"
)

func ParseAuthenticatedItem(ctx context.Context, wynnItem *entity.WynnItem, item *entity.AuthenticatedItem) (*entity.ItemInstance, error) {
	for id := range wynnItem.Stats {
		if _, ok := item.Stats[id]; !ok {
			return nil, fmt.Errorf("ID %s does not exist in item with Tracking Code\"%s\"", id, item.TrackingCode)
		}
	}

	return &entity.ItemInstance{
		Item:     item.Item,
		Stats:    item.Stats,
		WynnItem: wynnItem,
	}, nil
}

func ParseDecodedItem(ctx context.Context, decoded *decoder.DecodedItem, expected *entity.WynnItem) (*entity.ItemInstance, error) {

	canonicalIds := make(map[string]int, 0)
	for id, val := range decoded.Identifications {
		fmt.Printf("ID %d\n", id)
		name, ok := IdToName[int(id)]
		if !ok {
			return nil, fmt.Errorf("unknow identification numeric id: \"%d\"", id)
		}
		canonicalIds[name] = val
	}

	for id := range canonicalIds {
		if _, ok := expected.Stats[id]; !ok {
			fmt.Println(id)
			return nil, fmt.Errorf("ID %s does not exist in item \"%s\"", id, expected.Name)
		}
	}

	return &entity.ItemInstance{
		Item:     expected.Name,
		Stats:    canonicalIds,
		WynnItem: expected,
	}, nil

}

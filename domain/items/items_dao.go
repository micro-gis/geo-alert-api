package items

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/micro-gis/item-api/client/elasticsearch"
	"github.com/micro-gis/item-api/domain/queries"
	"github.com/micro-gis/utils/rest_errors"
	"strings"
)

const (
	indexItem = "item"
	typeItem  = "_doc"
)

func (i *Item) Save() rest_errors.RestErr {
	result, err := elasticsearch.Client.Index(indexItem, typeItem, i)
	if err != nil {
		return rest_errors.NewInternalServerError("error when trying to save item", errors.New("database error"))
	}
	i.Id = result.Id
	return nil
}

func (i *Item) Get() rest_errors.RestErr {
	itemId := i.Id
	result, err := elasticsearch.Client.Get(indexItem, typeItem, i.Id)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return rest_errors.NewNotFoundError(fmt.Sprintf("no items found with id : %s", i.Id))
		}
		return rest_errors.NewInternalServerError(fmt.Sprintf("error when trying get id %s", i.Id), errors.New("database error"))
	}
	bytes, err := result.Source.MarshalJSON()
	if err != nil {
		return rest_errors.NewInternalServerError("error when trying to parse database error", errors.New("database error"))
	}
	if err := json.Unmarshal(bytes, i); err != nil {
		return rest_errors.NewInternalServerError("error when trying to parse database error", errors.New("database error"))
	}
	i.Id = itemId
	return nil
}

func (i *Item) Search(query queries.EsQuery) ([]Item, rest_errors.RestErr) {
	result, err := elasticsearch.Client.Search(indexItem, query.Build())
	if err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to search documents", errors.New("database error"))
	}

	items := make([]Item, result.TotalHits())
	for i, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		item := Item{}
		if err := json.Unmarshal(bytes, &item); err != nil {
			return nil, rest_errors.NewInternalServerError("error when trying to parse response", errors.New("database error"))
		}
		item.Id = hit.Id
		items[i] = item
	}

	if len(items) == 0 {
		return nil, rest_errors.NewNotFoundError("not items found matching given criteria")
	}
	return items, nil
}

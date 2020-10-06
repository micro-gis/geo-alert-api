package items

import (
	"errors"
	"github.com/micro-gis/item-api/client/elasticsearch"
	"github.com/micro-gis/utils/rest_errors"
)

const (
	indexItem = "item"
)

func (i *Item) Save() *rest_errors.RestErr {
	result, err := elasticsearch.Client.Index(indexItem, i)
	if err != nil {
		return rest_errors.NewInternalServerError("error when trying to save item", errors.New("database error"))
	}
	i.Id = result.Id
	return nil
}

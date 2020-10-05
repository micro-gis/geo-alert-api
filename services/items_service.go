package services

import (
	"github.com/micro-gis/utils/rest_errors"
	"github.com/micro-gis/item-api/domain/items"
)

var (
	ItemService itemServiceInterface = &itemService{}
)

type itemServiceInterface interface {
	Create(items.Item) (*items.Item, *rest_errors.RestErr)
	Get(string) (*items.Item, *rest_errors.RestErr)
}

type itemService struct{}

func (i itemService) Create(item items.Item) (*items.Item, *rest_errors.RestErr) {
	return nil, rest_errors.NewBadRequestError("implement me")
}

func (i itemService) Get(s string) (*items.Item, *rest_errors.RestErr) {
	return nil, rest_errors.NewBadRequestError("implement me")
}

package controllers

import (
	"fmt"
	"github.com/micro-gis/oauth-go/oauth"
	"github.com/micro-gis/item-api/domain/items"
	"net/http"
)

func Create(w http.ResponseWriter, r *http.Request) {
	if err := oauth.AuthenticateRequest(r); err != nil {
		//TODO: Return error to the caller
		return
	}

	item := items.Item{
		Id:                "",
		Seller:            oauth.GetCallerId(r),
		Title:             "",
		Description:       items.Description{},
		Pictures:          nil,
		Video:             "",
		Price:             0,
		AvailableQuantity: 0,
		SoldQuantity:      0,
		Status:            "",
	}
	result, err := services.ItemsService.Create(item)
	if err != nil {
		//TODO: Return error json to the user
	}
	fmt.Println(result)
	//TODO : Return created item with HTTP status 201 -Created
	
}

func Get(w http.ResponseWriter, r *http.Request) {

}

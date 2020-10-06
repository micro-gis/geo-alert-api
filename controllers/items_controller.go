package controllers

import (
	"encoding/json"
	"github.com/micro-gis/item-api/domain/items"
	"github.com/micro-gis/item-api/services"
	"github.com/micro-gis/item-api/utils/http_utils"
	"github.com/micro-gis/oauth-go/oauth"
	"github.com/micro-gis/utils/rest_errors"
	"io/ioutil"
	"net/http"
)

var (
	ItemController itemsControllerInterface = &itemsController{}
)

type itemsControllerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}

type itemsController struct {
}

func (c *itemsController) Create(w http.ResponseWriter, r *http.Request) {
	if err := oauth.AuthenticateRequest(r); err != nil {
		http_utils.ResponseError(w, err)
		return
	}

	sellerId := oauth.GetCallerId(r)

	if sellerId == 0 {
		restErr := rest_errors.NewUnauthorizedError()
		http_utils.ResponseError(w, restErr)
		return
	}
	var itemRequest items.Item
	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		respErr := rest_errors.NewBadRequestError("invalid request body")
		http_utils.ResponseError(w, respErr)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(requestBody, &itemRequest); err != nil {
		respErr := rest_errors.NewBadRequestError("invalid request body")
		http_utils.ResponseError(w, respErr)
		return
	}

	itemRequest.Seller = sellerId
	result, createErr := services.ItemService.Create(itemRequest)
	if createErr != nil {
		http_utils.ResponseError(w, createErr)
		return
	}
	http_utils.ResponseJson(w, http.StatusCreated, result)
	return

}

func (c *itemsController) Get(w http.ResponseWriter, r *http.Request) {

}

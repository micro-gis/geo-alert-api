package app

import (
	"github.com/micro-gis/item-api/controllers"
	"net/http"
)

func mapUrls() {
	router.HandleFunc("/items", controllers.ItemController.Create).Methods(http.MethodPost)
	router.HandleFunc("/ping", controllers.PingController.Ping).Methods(http.MethodGet)
	router.HandleFunc("/items/{id}", controllers.ItemController.Get).Methods(http.MethodGet)
	router.HandleFunc("/items/search", controllers.ItemController.Search).Methods(http.MethodPost)
}

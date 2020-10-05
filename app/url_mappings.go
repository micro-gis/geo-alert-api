package app

import (
	"github.com/micro-gis/item-api/controllers"
	"net/http"
)

func mapUrls() {
	router.HandleFunc("/items", controllers.ItemController.Create).Methods(http.MethodPost)
	router.HandleFunc("/ping", controllers.PingController.Ping).Methods(http.MethodGet)
}

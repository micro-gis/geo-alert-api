package app

import (
	"github.com/micro-gis/geo-alert-api/controllers"
	"net/http"
)

func mapUrls() {
	router.HandleFunc("/geoalerts", controllers.GeoAlertController.Create).Methods(http.MethodPost)
	router.HandleFunc("/ping", controllers.PingController.Ping).Methods(http.MethodGet)
	router.HandleFunc("/geoalerts/{id}", controllers.GeoAlertController.Get).Methods(http.MethodGet)
	router.HandleFunc("/geoalerts/search", controllers.GeoAlertController.Search).Methods(http.MethodPost)
	router.HandleFunc("/geoalerts/user/{user_id}", controllers.GeoAlertController.GetUserGeoAlerts).Methods(http.MethodGet)
}

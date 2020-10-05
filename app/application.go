package app

import (
	"github.com/gorilla/mux"
	"github.com/micro-gis/item-api/client/elasticsearch"
	"net/http"
	"time"
)

var (
	router = mux.NewRouter()
)

func StartApplication() {
	elasticsearch.Init()
	mapUrls()
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 500 * time.Millisecond,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  2 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

package router

import (
	"encoding/json"
	"net/http"
	"server/middleware"
	bookRoutes "server/router/book"

	"github.com/gorilla/mux"
)

func healthCheck(res http.ResponseWriter, req *http.Request) {
	json.NewEncoder(res).Encode("Hey!")
}

func Routes() *mux.Router {
	router := mux.NewRouter()

	// set middleware
	router.Use(middleware.HeaderMiddleware)
	router.Use(middleware.LogRequest)

	router.HandleFunc("/", healthCheck).Methods("GET", "OPTIONS")

	// set subroutes
	book := router.PathPrefix("/books").Subrouter()
	bookRoutes.Router(book)

	return router
}

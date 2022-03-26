package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/middleware"
	bookRoutes "server/router/book"

	"github.com/gorilla/mux"
)

func healthCheck(res http.ResponseWriter, req *http.Request) {
	json.NewEncoder(res).Encode("Welcome to the book store!")
}

func Routes() *mux.Router {
	router := mux.NewRouter()

	// set middleware
	router.Use(middleware.HeaderMiddleware)
	router.Use(middleware.LogRequest)

	// set subroutes
	book := router.PathPrefix("/books").Subrouter()
	bookRoutes.Router(book)

	router.HandleFunc("/", healthCheck).Methods("GET", "OPTIONS")
	fmt.Println("server started")
	return router
}

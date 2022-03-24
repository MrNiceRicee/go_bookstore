package bookRoutes

import (
	"github.com/gorilla/mux"
)

func Router(router *mux.Router) {
	router.Path("").HandlerFunc(Search).Queries("limit", "{id:[0-9]+}").Methods("GET")
	router.Path("").HandlerFunc(Search).Methods("GET")
	router.Path("").HandlerFunc(Create).Methods("POST")
	router.Path("/{id}").HandlerFunc(Update).Methods("PUT")
	router.Path("/{id}").HandlerFunc(Delete).Methods("DELETE")
}

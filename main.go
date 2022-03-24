package main

import (
	"net/http"
	"server/connection"
	"server/router"
)

func main() {
	app := router.Routes()
	connection.CreateConnection()

	http.ListenAndServe("localhost:8000", app)
}

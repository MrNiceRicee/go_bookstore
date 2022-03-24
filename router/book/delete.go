package bookRoutes

import (
	"errors"
	"net/http"
	"server/connection"
	response "server/router/handlers"
	"strconv"

	"github.com/gorilla/mux"
)

func deleteBook(id int) error {
	db := connection.DB

	query := `DELETE FROM "Books" WHERE "_id"=$1`

	_, err := db.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func validateDelete(req *http.Request) (int, error) {
	// grab param
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return 0, errors.New("missing Id")
	}
	_, err = findBook(id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func Delete(res http.ResponseWriter, req *http.Request) {

	id, err := validateDelete(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = deleteBook(id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	response.HandleResponse(res, "", http.StatusNoContent)
}

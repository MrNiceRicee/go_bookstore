package bookRoutes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/connection"
	"server/models"
	response "server/router/handlers"
)

func validateCreate(req *http.Request) (models.Book, error) {
	var book models.Book
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&book)

	if err != nil {
		return book, errors.New("invalid request body")
	}
	if book.Author == "" {
		return book, errors.New("missing author")
	}
	if book.Title == "" {
		return book, errors.New("missing title")
	}

	return book, nil
}

func createBook(book models.Book) (models.Book, error) {
	db := connection.DB

	query := `INSERT INTO "Books"("author", "title", "genre")
						VALUES($1, $2, $3)
						RETURNING "_id"
						`

	var created models.Book

	err := db.QueryRow(query, book.Author, book.Title, book.Genre).Scan(&created.Id)

	if err != nil {
		return created, err
	}

	return created, nil
}

func Create(res http.ResponseWriter, req *http.Request) {

	book, err := validateCreate(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	results, err := createBook(book)

	if err != nil {
		http.Error(res, "Internal Error", http.StatusInternalServerError)
		return

	}

	response.HandleResponse(res, fmt.Sprintf("Created Book. %v", results.Id), http.StatusOK)
}

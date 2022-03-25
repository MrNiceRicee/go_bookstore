package bookRoutes

import (
	"encoding/json"
	"errors"
	"net/http"
	"server/connection"
	"server/models"
	response "server/router/handlers"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type validBook struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Genre  string `json:"genre"`
}

func findBook(id int) (models.Book, error) {
	db := connection.DB
	var book models.Book

	query := `SELECT "_id","title", "author", "genre", "createdAt", "updatedAt"
						FROM "Books"
						WHERE "_id"=$1
						`

	row := db.QueryRow(query, id)
	err := row.Scan(&book.Id,
		&book.Title,
		&book.Author,
		&book.Genre,
		&book.CreatedAt,
		&book.UpdatedAt)

	if err != nil {
		return book, errors.New("book not found")
	}

	// return book from db
	return book, nil
}

func validateUpdatePackage(req *http.Request) (models.Book, error) {
	var book models.Book
	var checkBook validBook
	// grab param
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return book, errors.New("missing Id")
	}

	// build book from DB
	book, err = findBook(id)
	if err != nil {
		return book, err
	}

	// decode body & grab valid fields only
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&checkBook)
	if err != nil {
		return book, errors.New("invalid request body")
	}

	// probably not the best way?
	if checkBook.Author != "" {
		book.Author = checkBook.Author
	}
	if checkBook.Title != "" {
		book.Title = checkBook.Title
	}
	if checkBook.Genre != "" {
		book.Genre = checkBook.Genre
	}
	return book, nil
}

func updateBook(book models.Book) (models.Book, error) {
	db := connection.DB

	query := `UPDATE "Books"
						SET "title"=$1,
								"author"=$2,
								"genre"=$3,
								"updatedAt"=$4
						WHERE "_id"=$5
						RETURNING "_id",
											"title",
											"author",
											"genre",
											"createdAt",
											"updatedAt"
						`

	var update models.Book
	row := db.QueryRow(query,
		book.Title,
		book.Author,
		book.Genre,
		time.Now(),
		book.Id)

	err := row.Scan(&update.Id,
		&update.Title,
		&update.Author,
		&update.Genre,
		&update.CreatedAt,
		&update.UpdatedAt)
	if err != nil {
		return update, err
	}

	return update, nil

}

func Update(res http.ResponseWriter, req *http.Request) {

	update, err := validateUpdatePackage(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	results, err := updateBook(update)

	if err != nil {
		http.Error(res, "Internal Error", http.StatusInternalServerError)
		return
	}

	response.HandleResponse(res, results, http.StatusOK)
}

package bookRoutes

import (
	"errors"
	"net/http"
	"server/connection"
	"server/models"
	response "server/router/handlers"
	"strconv"
)

type validQueries struct {
	Limit int `json:"limit"`
}

type searchResults struct {
	Data   []models.Book `json:"data"`
	Length int           `json:"length"`
}

func defaultQuery() validQueries {
	return validQueries{
		Limit: 25,
	}
}

func searchDB(queryParams validQueries) ([]models.Book, error) {
	db := connection.DB

	var books []models.Book

	query := `SELECT "_id", "title", "author", "genre", "createdAt", "updatedAt"
						FROM "Books"
						LIMIT $1`

	rows, err := db.Query(query, queryParams.Limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var book models.Book
		err = rows.Scan(&book.Id,
			&book.Title,
			&book.Author,
			&book.Genre,
			&book.CreatedAt,
			&book.UpdatedAt)

		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}
	return books, nil
}

func validateQueries(req *http.Request) (validQueries, error) {
	baseQuery := defaultQuery()

	// only return error if there is an entered limit
	receivedLimit := req.FormValue("limit")
	if receivedLimit != "" {
		limit, err := strconv.Atoi(receivedLimit)
		if err != nil {
			return baseQuery, errors.New("limit must be a number")
		}
		baseQuery.Limit = limit
	}
	return baseQuery, nil
}

func Search(res http.ResponseWriter, req *http.Request) {

	// not sure about the validation for queries
	query, err := validateQueries(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := searchDB(query)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	response.HandleCustomResponse(res, searchResults{
		Data:   results,
		Length: len(results),
	}, http.StatusOK)
}

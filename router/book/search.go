package bookRoutes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/connection"
	"server/models"
	response "server/router/handlers"
	"server/router/utility"

	"strconv"

	sq "github.com/Masterminds/squirrel"
)

type validQueries struct {
	Limit  int           `json:"limit,omitempty"`
	Filter models.Filter `json:"filter,omitempty"`
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

func getValidFilters() []string {
	return []string{
		"_id", "author", "title", "genre",
	}
}

func searchDB(queryParams validQueries) ([]models.Book, error) {
	db := connection.DB

	var books []models.Book

	selectBooks := sq.Select("_id", "title", "author", "genre", `"createdAt"`, `"updatedAt"`).From(`"Books"`)

	selectBooks = utility.BuildWhere(selectBooks, getValidFilters(), queryParams.Filter)
	selectBooks = selectBooks.Limit(uint64(queryParams.Limit)).PlaceholderFormat(sq.Dollar)

	sql, args, err := selectBooks.ToSql()
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s, %+v\n", sql, args)
	rows, err := db.Query(sql, args...)

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

	var receivedFilter models.Filter
	err := json.Unmarshal([]byte(req.FormValue("filter")), &receivedFilter)
	if err != nil {
		// filter is not required, so we can ignore the error
		return baseQuery, nil
	}
	baseQuery.Filter = receivedFilter

	return baseQuery, nil
}

func Search(res http.ResponseWriter, req *http.Request) {

	// not sure about the validation for queries
	query, err := validateQueries(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
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

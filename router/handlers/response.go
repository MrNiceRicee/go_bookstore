package response

import (
	"encoding/json"
	"net/http"
	"server/models"
)

func HandleResponse(res http.ResponseWriter, data interface{}, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	res.WriteHeader(statusCode)
	json.NewEncoder(res).Encode(models.Response{
		Data: data,
	})
}

func HandleCustomResponse(res http.ResponseWriter, data interface{}, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	res.WriteHeader(statusCode)
	json.NewEncoder(res).Encode(data)
}

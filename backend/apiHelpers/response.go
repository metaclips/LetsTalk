package apihelpers

import (
	"encoding/json"
	"net/http"
)

// ResponseData structure
type response struct {
	Data       interface{} `json:"data"`
	Error      string      `json:"error"`
	StatusCode int         `json:"statusCode"`
}

// Respond returns basic response structure
func Respond(w http.ResponseWriter, data interface{}, statusCode int, err error) error {
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(response{
		Data:       data,
		Error:      err.Error(),
		StatusCode: statusCode,
	})
}

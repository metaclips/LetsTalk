package apihelpers

import (
	"encoding/json"
	"net/http"
)

// ResponseData structure
type response struct {
	Data       interface{} `json:"data"`
	Msg        string      `json:"msg"`
	Error      string      `json:"error"`
	StatusCode int         `json:"statusCode"`
}

// Respond returns basic response structure
func Respond(w http.ResponseWriter, statusCode int, data interface{}, msg, err string) error {
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(response{
		Data:       data,
		Error:      err,
		Msg:        msg,
		StatusCode: statusCode,
	})
}

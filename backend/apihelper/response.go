package apihelper

import (
	"encoding/json"
	"net/http"
)

// ResponseData structure
type Response struct {
	Data       interface{} `json:"data"`
	Msg        string      `json:"msg"`
	Error      string      `json:"error"`
	StatusCode int         `json:"statusCode"`
}

// Respond returns basic response structure
func Respond(w http.ResponseWriter, statusCode int, data interface{}, msg, errVal string) error {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")

	resp := Response{
		Data:       data,
		Error:      errVal,
		Msg:        msg,
		StatusCode: statusCode,
	}

	return json.NewEncoder(w).Encode(resp)
}

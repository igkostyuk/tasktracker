package web

import (
	"encoding/json"
	"net/http"
)

// Respond converts a Go value to JSON and sends it to the client.
func Respond(w http.ResponseWriter, data interface{}, statusCode int) {

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	w.Write(jsonData)
}

// RespondError example
func RespondError(w http.ResponseWriter, err error, status int) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	Respond(w, er, status)
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

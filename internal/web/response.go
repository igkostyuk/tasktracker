package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/igkostyuk/tasktracker/internal/middleware"
)

// Respond converts a Go value to JSON and sends it to the client.
func Respond(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		RespondError(w, r, err, http.StatusInternalServerError)

		return
	}
	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")
	// Write the status code to the response.
	w.WriteHeader(statusCode)
	// Send the result back to the client.
	_, err = w.Write(jsonData)
	logError(r, err)
}

// RespondError sent error JSON message to the client.
func RespondError(w http.ResponseWriter, r *http.Request, err error, status int) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	if status == http.StatusInternalServerError {
		logError(r, err)
	}
	Respond(w, r, er, status)
}

// HTTPError represent error respons message.
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func logError(r *http.Request, err error) {
	logEntry := middleware.GetLogEntry(r)
	if logEntry != nil {
		logEntry.WriteError(err)
	} else {
		log.Printf("%s: ERROR:\n%s", middleware.GetRequestID(r.Context()), err.Error())
	}
}

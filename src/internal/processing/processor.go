package processing

import (
	"encoding/json"
	"net/http"

	"github.com/emurray647/audioServer/internal/dbconnector"
	"github.com/emurray647/audioServer/internal/model"
)

// RequestProcessor is a struct to handle the processor of all the requests
type RequestProcessor struct {
	db         *dbconnector.DBConnection
	filePrefix string
}

// Creates a RequestProcesesor object
// db - DBConnection object to communicate with the DB
// filePrefix - the location on disk to store our files
func NewRequestProcessor(db *dbconnector.DBConnection, filePrefix string) *RequestProcessor {
	return &RequestProcessor{
		db:         db,
		filePrefix: filePrefix,
	}
}

// Handles requests for an undefined endpoint
func (RequestProcessor) UnknownPath(w http.ResponseWriter, r *http.Request) {
	setStatus(w, http.StatusNotFound, "path does not exist")
}

// set the error status on the ResponseWriter
func setStatus(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	sm := model.StatusMessage{
		StatusCode: statusCode,
		Message:    message,
		Success:    false,
	}

	json.NewEncoder(w).Encode(sm)
}

// set the ResponseWriter with a success message
func setSuccess(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	sm := model.StatusMessage{
		StatusCode: statusCode,
		Message:    message,
		Success:    true,
	}

	json.NewEncoder(w).Encode(sm)
}

// func setErrorStatus(w http.ResponseWriter, statusCode int, err error) {
// 	w.WriteHeader(statusCode)
// 	w.Header().Set("Content-Type", "application/json")

// 	sm := model.StatusMessage{
// 		StatusCode: statusCode,
// 		Message:    err.Error(),
// 		Success:    false,
// 	}

// 	json.NewEncoder(w).Encode(sm)
// }

package processing

import (
	"net/http"

	"github.com/emurray647/audioServer/internal/dbconnector"
)

type RequestProcessor struct {
	db         *dbconnector.DBConnection
	filePrefix string
}

func NewRequestProcessor(db *dbconnector.DBConnection, filePrefix string) *RequestProcessor {
	return &RequestProcessor{
		db:         db,
		filePrefix: filePrefix,
	}
}

func (RequestProcessor) UnknownPath(w http.ResponseWriter, r *http.Request) {
	setStatus(w, http.StatusNotFound, "path does not exist", false)
}

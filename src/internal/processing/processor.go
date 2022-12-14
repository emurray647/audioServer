package processing

import (
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

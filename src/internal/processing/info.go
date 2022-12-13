package processing

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/emurray647/audioServer/internal/dbconnector"
	"github.com/emurray647/audioServer/internal/model"
	log "github.com/sirupsen/logrus"
)

func Info(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	if name == "" {
		log.Info("Info did not receive name param")
		setStatus(w, http.StatusBadRequest, "did not receive name param", false)
		return
	}

	details, err := info(name)
	if err != nil && errors.Is(err, fileDoesNotExist) {
		setStatus(w, http.StatusNotFound, fileDoesNotExist.Error(), false)
		return
	} else if err != nil {
		setStatus(w, http.StatusInternalServerError, "unknown error", false)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)

}

func info(name string) (*model.WavFileDetails, error) {
	dbConnection, err := dbconnector.OpenDBConnection()
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %w", err)

	}
	defer dbConnection.Close()

	details, err := dbConnection.GetWavDetails(name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fileDoesNotExist
		}
		return nil, fmt.Errorf("failed getting wav details: %w", err)
	}

	return details, nil
}

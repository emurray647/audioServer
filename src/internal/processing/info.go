package processing

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/emurray647/audioServer/internal/model"
	log "github.com/sirupsen/logrus"
)

// This file holds the logic for handling the /info endpoint

func (p *RequestProcessor) Info(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	if name == "" {
		log.Info("Info did not receive name param")
		setStatus(w, http.StatusBadRequest, "did not receive name param")
		return
	}

	details, err := p.info(name)
	if err != nil && errors.Is(err, fileDoesNotExist) {
		setStatus(w, http.StatusNotFound, fileDoesNotExist.Error())
		return
	} else if err != nil {
		setStatus(w, http.StatusInternalServerError, "unknown error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)

}

func (p *RequestProcessor) info(name string) (*model.AudioFileDetails, error) {
	details, err := p.db.GetFileDetails(name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fileDoesNotExist
		}
		return nil, fmt.Errorf("failed getting file details: %w", err)
	}

	return details, nil
}

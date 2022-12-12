package processing

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/emurray647/audioServer/internal/model"
	log "github.com/sirupsen/logrus"
)

// type customError struct {
// 	s string
// }

// func createCustomError(msg string) error {
// 	return &customError{
// 		s: msg,
// 	}
// }

// func (e *customError) Error() string {
// 	return e.s
// }

var (
	invalidWAVError error = errors.New("invalid WAV file")

	invalidFileFormat error = errors.New("invalid file format")
	unknownFileType   error = errors.New("unknown file extension")

	fileAlreadyExists error = errors.New("file already exists")

	fileDoesNotExist error = errors.New("file does not exist")
)

func logError(w http.ResponseWriter, statusCode int, err error) {

	log.Errorf(err.Error())

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	var msg string
	if errors.Is(err, invalidWAVError) {
		msg = invalidWAVError.Error()
	} else if errors.Is(err, fileAlreadyExists) {
		msg = fileAlreadyExists.Error()
	} else if errors.Is(err, fileDoesNotExist) {
		msg = fileDoesNotExist.Error()
	} else {
		msg = "unknown error"
	}

	sm := model.StatusMessage{
		StatusCode: statusCode,
		Message:    msg,
		Success:    false,
	}

	err = json.NewEncoder(w).Encode(sm)
	if err != nil {
		panic(err)
	}
}

package processing

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/emurray647/audioServer/internal/dbconnector"
	"github.com/emurray647/audioServer/internal/format"
	"github.com/emurray647/audioServer/internal/model"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	writePrefix = "/data"
)

func Upload(w http.ResponseWriter, r *http.Request) {

	fmt.Println("upload")

	// grab the data from the request and grab the name parameter (if applicable)
	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("could not read request body: %w", err)
		setStatus(w, http.StatusBadRequest, "could not read POST body", false)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		name = generateName(buffer)
	}

	// upload the file
	err = upload(name, buffer)
	if err != nil {
		log.Errorf("error uploading file: %w", err)
		if errors.Is(err, fileAlreadyExists) {
			setStatus(w, http.StatusConflict, fileAlreadyExists.Error(), false)
		} else if errors.Is(err, invalidWAVError) {
			setStatus(w, http.StatusBadRequest, invalidWAVError.Error(), false)
		} else {
			setStatus(w, http.StatusInternalServerError, "unknown error", false)
		}
		return
	}

}

func Delete(w http.ResponseWriter, r *http.Request) {

	fmt.Println("delete")

	// grab off the filename variable
	vars := mux.Vars(r)
	filename, ok := vars["filename"]
	if !ok {
		log.Errorf("no file provided to delete")
		http.Error(w, "did not provide file to delete", http.StatusBadRequest)
		return
	}

	err := delete(filename)
	if err != nil {
		log.Errorf("could not delete file: %w", err)
		if errors.Is(err, fileDoesNotExist) {
			setStatus(w, http.StatusBadRequest, fileDoesNotExist.Error(), false)
		} else {
			setStatus(w, http.StatusInternalServerError, "unknown error", false)
		}
		return
	}
}

func Download(w http.ResponseWriter, r *http.Request) {

	fmt.Println("download")

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "no name value provided", http.StatusBadRequest)
		return
	}

	buffer, err := download(name)
	if err != nil {
		http.Error(w, fmt.Sprintf("error retrieving file: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(buffer)
}

func upload(filename string, data []byte) error {
	// open db connection
	dbConnection, err := dbconnector.OpenDBConnection()
	if err != nil {
		return fmt.Errorf("could not open database connection: %w", err)

	}
	defer dbConnection.Close()

	// make sure this is not a duplicate entry
	count, err := dbConnection.CountWavFiles(filename)
	if err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}
	if count > 0 {
		return fileAlreadyExists
	}

	// before we write this file, we should verify it is a wav
	// as well as get some stats about it
	details, err := format.ParseFile(filename, data)
	if err != nil && errors.Is(err, &format.InvalidFileError{}) {
		return invalidWAVError
	}

	// write the file to disk
	fullpath := fmt.Sprintf("%s/%s", writePrefix, filename)
	err = os.WriteFile(fullpath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file to disk: %w", err)
	}

	// write the info to the DB
	wav := &model.WavFile{
		WavFileDetails: *details,
		URI:            fullpath,
	}
	err = dbConnection.AddWavFile(wav)
	if err != nil {
		return fmt.Errorf("failed to put wav in db: %w", err)
	}

	return nil
}

func delete(filename string) error {
	dbConnection, err := dbconnector.OpenDBConnection()
	if err != nil {
		return fmt.Errorf("could not open database connection: %w", err)

	}
	defer dbConnection.Close()

	fileURI, err := dbConnection.DeleteWav(filename)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return fileDoesNotExist
	} else if err != nil {
		return fmt.Errorf("could not delete file: %w", err)
	}

	if err = os.Remove(fileURI); err != nil {
		return fmt.Errorf("could not remove file %s: %w", fileURI, err)
	}

	return nil
}

func download(filename string) ([]byte, error) {
	dbConnection, err := dbconnector.OpenDBConnection()
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %w", err)

	}
	defer dbConnection.Close()

	uri, err := dbConnection.GetWavURI(filename)
	if err != nil {
		return nil, fmt.Errorf("could not find wav in db: %w", err)
	}

	fileBytes, err := ioutil.ReadFile(uri)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	return fileBytes, nil
}

func generateName(buffer []byte) string {
	hash := md5.Sum(buffer)
	return fmt.Sprintf("%s.wav", hex.EncodeToString(hash[:]))
}

func setStatus(w http.ResponseWriter, statusCode int, message string, success bool) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	sm := model.StatusMessage{
		StatusCode: statusCode,
		Message:    message,
		Success:    success,
	}

	json.NewEncoder(w).Encode(sm)
}

func setErrorStatus(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	sm := model.StatusMessage{
		StatusCode: statusCode,
		Message:    err.Error(),
		Success:    false,
	}

	json.NewEncoder(w).Encode(sm)
}

package processing

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/emurray647/audioServer/internal/dbconnector"
	"github.com/emurray647/audioServer/internal/model"
	"github.com/gorilla/mux"
)

const (
	writePrefix = "/data"
)

func Upload(w http.ResponseWriter, r *http.Request) {

	// grab the data from the request and grab the name parameter (if applicable)
	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read request body: %v", err), http.StatusBadRequest)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		name = generateName(buffer)
	}

	// upload the file
	err = upload(name, buffer)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not open database connection: %v", err), http.StatusInternalServerError)
		return
	}

}

func Delete(w http.ResponseWriter, r *http.Request) {
	// grab off the filename variable
	vars := mux.Vars(r)
	filename, ok := vars["filename"]
	if !ok {
		http.Error(w, "did not provide file to delete", http.StatusBadRequest)
		return
	}

	err := delete(filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not delete file: %v", err), http.StatusBadRequest)
	}
}

func Download(w http.ResponseWriter, r *http.Request) {
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
		return fmt.Errorf("entry already exists")
	}

	// before we write this file, we should verify it is a wav
	// as well as get some stats about it
	details, err := parseWav(filename, data)
	if err != nil && errors.Is(err, &invalidFileError{}) {
		return fmt.Errorf("invalid wav file: %w", err)
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

	err = dbConnection.DeleteWav(filename)
	if err != nil {
		return fmt.Errorf("could not delete file: %w", err)
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

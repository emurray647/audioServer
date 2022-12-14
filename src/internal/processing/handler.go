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

	"github.com/emurray647/audioServer/internal/format"
	"github.com/emurray647/audioServer/internal/model"
	log "github.com/sirupsen/logrus"
)

func (p *RequestProcessor) Upload(w http.ResponseWriter, r *http.Request) {

	// grab the data from the request and grab the name parameter (if applicable)
	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("could not read request body: %w", err)
		setStatus(w, http.StatusInternalServerError, "could not read POST body", false)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		name = generateName(buffer)
	}

	// upload the file
	err = p.upload(name, buffer)
	if err != nil {
		err = fmt.Errorf("error uploading file: %w", err)
		log.Errorf(err.Error())
		if errors.Is(err, fileAlreadyExists) {
			setStatus(w, http.StatusConflict, fileAlreadyExists.Error(), false)
		} else if errors.Is(err, invalidFileFormat) {
			setStatus(w, http.StatusBadRequest, invalidFileFormat.Error(), false)
		} else if errors.Is(err, unknownFileType) {
			setStatus(w, http.StatusBadRequest, unknownFileType.Error(), false)
		} else {
			setStatus(w, http.StatusInternalServerError, "unknown error", false)
		}
		return
	}

	setSuccess(w, http.StatusOK, "Successfully uploaded file")

}

func (p *RequestProcessor) Delete(w http.ResponseWriter, r *http.Request) {

	fmt.Println("delete")

	// grab off the filename variable
	filename := r.URL.Query().Get("name")
	if filename == "" {
		setStatus(w, http.StatusBadRequest, "did not provide file to delete", false)
		return
	}

	err := p.delete(filename)
	if err != nil {
		log.Errorf("could not delete file: %w", err)
		if errors.Is(err, fileDoesNotExist) {
			setStatus(w, http.StatusBadRequest, fileDoesNotExist.Error(), false)
		} else {
			setStatus(w, http.StatusInternalServerError, "unknown error", false)
		}
		return
	}

	setSuccess(w, http.StatusOK, fmt.Sprintf("Successfully deleted file %s", filename))
}

func (p *RequestProcessor) Download(w http.ResponseWriter, r *http.Request) {

	fmt.Println("download")

	name := r.URL.Query().Get("name")
	if name == "" {
		setStatus(w, http.StatusBadRequest, "no name value provided", false)
		return
	}

	buffer, err := p.download(name)
	if err != nil && errors.Is(err, fileDoesNotExist) {
		setStatus(w, http.StatusNotFound, fileDoesNotExist.Error(), false)
		return
	} else if err != nil {
		log.Errorf(err.Error())
		setStatus(w, http.StatusInternalServerError, "unknown error", false)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(buffer)

}

func (p *RequestProcessor) upload(filename string, data []byte) error {
	// make sure this is not a duplicate entry
	count, err := p.db.CountWavFiles(filename)
	if err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}
	if count > 0 {
		return fileAlreadyExists
	}

	// before we write this file, we should verify it is a valid file
	// as well as get some stats about it
	details, err := format.ParseFile(filename, data)
	if err != nil && errors.Is(err, format.InvalidFile) {
		return invalidFileFormat
	} else if err != nil && errors.Is(err, format.UnknownFormat) {
		return unknownFileType
	} else if err != nil {
		return err
	}

	// write the file to disk
	fullpath := fmt.Sprintf("%s/%s", p.filePrefix, details.Name)
	err = os.WriteFile(fullpath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file to disk: %w", err)
	}

	// write the info to the DB
	wav := &model.WavFile{
		WavFileDetails: *details,
		URI:            fullpath,
	}
	err = p.db.AddWavFile(wav)
	if err != nil {
		return fmt.Errorf("failed to put wav in db: %w", err)
	}

	return nil
}

func (p *RequestProcessor) delete(filename string) error {
	fileURI, err := p.db.DeleteWav(filename)
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

func (p *RequestProcessor) download(filename string) ([]byte, error) {
	uri, err := p.db.GetWavURI(filename)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fileDoesNotExist
		}
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
	return fmt.Sprintf("%s", hex.EncodeToString(hash[:]))
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

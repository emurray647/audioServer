package processing

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/emurray647/audioServer/internal/format"
	"github.com/emurray647/audioServer/internal/model"
	log "github.com/sirupsen/logrus"
)

// This file contains all the logic for the /files and /download endpoints
// (basically the Create/Read/Delete endpiosn)

func (p *RequestProcessor) Upload(w http.ResponseWriter, r *http.Request) {
	// grab the data from the request and grab the name parameter (if applicable)
	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("could not read request body: %w", err)
		setStatus(w, http.StatusInternalServerError, "could not read POST body")
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
			setStatus(w, http.StatusConflict, fileAlreadyExists.Error())
		} else if errors.Is(err, invalidFileFormat) {
			setStatus(w, http.StatusBadRequest, invalidFileFormat.Error())
		} else if errors.Is(err, unknownFileType) {
			setStatus(w, http.StatusBadRequest, unknownFileType.Error())
		} else {
			setStatus(w, http.StatusInternalServerError, "unknown error")
		}
		return
	}

	setSuccess(w, http.StatusOK, "Successfully uploaded file")

}

func (p *RequestProcessor) Delete(w http.ResponseWriter, r *http.Request) {

	// grab off the filename variable
	filename := r.URL.Query().Get("name")
	if filename == "" {
		setStatus(w, http.StatusBadRequest, "did not provide file to delete")
		return
	}

	err := p.delete(filename)
	if err != nil {
		log.Errorf("could not delete file: %w", err)
		if errors.Is(err, fileDoesNotExist) {
			setStatus(w, http.StatusBadRequest, fileDoesNotExist.Error())
		} else {
			setStatus(w, http.StatusInternalServerError, "unknown error")
		}
		return
	}

	setSuccess(w, http.StatusOK, fmt.Sprintf("Successfully deleted file %s", filename))
}

func (p *RequestProcessor) Download(w http.ResponseWriter, r *http.Request) {
	// grab the name param
	name := r.URL.Query().Get("name")
	if name == "" {
		setStatus(w, http.StatusBadRequest, "no name value provided")
		return
	}

	buffer, err := p.download(name)
	if err != nil && errors.Is(err, fileDoesNotExist) {
		setStatus(w, http.StatusNotFound, fileDoesNotExist.Error())
		return
	} else if err != nil {
		log.Errorf(err.Error())
		setStatus(w, http.StatusInternalServerError, "unknown error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(buffer)

}

func (p *RequestProcessor) upload(filename string, data []byte) error {
	// make sure this is not a duplicate entry
	count, err := p.db.CountFiles(filename)
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
	audioFile := &model.AudioFile{
		AudioFileDetails: *details,
		URI:              fullpath,
	}
	err = p.db.AddFile(audioFile)
	if err != nil {
		return fmt.Errorf("failed to put file in db: %w", err)
	}

	return nil
}

func (p *RequestProcessor) delete(filename string) error {
	// first delete the file from the database
	fileURI, err := p.db.DeleteFile(filename)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return fileDoesNotExist
	} else if err != nil {
		return fmt.Errorf("could not delete file: %w", err)
	}

	// if the line was successfully removed from the DB, then we delete the actual file
	if err = os.Remove(fileURI); err != nil {
		return fmt.Errorf("could not remove file %s: %w", fileURI, err)
	}
	return nil
}

func (p *RequestProcessor) download(filename string) ([]byte, error) {
	// get the uri of the file
	uri, err := p.db.GetFileURI(filename)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fileDoesNotExist
		}
		return nil, fmt.Errorf("could not find file in db: %w", err)
	}

	// copy the file into a buffer to return to the caller
	fileBytes, err := ioutil.ReadFile(uri)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	return fileBytes, nil
}

// generate a name based on the buffer comments
func generateName(buffer []byte) string {
	hash := md5.Sum(buffer)
	return fmt.Sprintf("%s", hex.EncodeToString(hash[:]))
}

package processing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/emurray647/audioServer/internal/dbconnector"
	"github.com/emurray647/audioServer/internal/model"
)

func List(w http.ResponseWriter, r *http.Request) {
	// get all elements that match query

	values := r.URL.Query()

	details, err := list(values)
	if err != nil {
		// log.Errorf("failed listing wav details: %w", err)
		// setStatus(w, http.StatusInternalServerError, "unknown error", false)
		logError(w, http.StatusBadRequest, err)
		return
	}

	err = json.NewEncoder(w).Encode(details)
	if err != nil {
		// log.Errorf("failed to encode wav JSON")
		// setStatus(w, http.StatusInternalServerError, "unknown error", false)
		logError(w, http.StatusInternalServerError, fmt.Errorf("failed to encode wav JSON: %w", err))
		return
	}

	// // var wavs model.WavFiles
	// wavs, err := getData()
	// if err != nil {
	// 	http.Error(w, "failed to read from database", http.StatusInternalServerError)
	// 	return
	// }

	// err = json.NewEncoder(w).Encode(&wavs)
	// if err != nil {
	// 	http.Error(w, "failed to encode wav data", http.StatusInternalServerError)
	// 	return
	// }

}

func list(values url.Values) (*model.WavFilesDetailsSlice, error) {
	dbConnection, err := dbconnector.OpenDBConnection()
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %w", err)

	}
	defer dbConnection.Close()

	result, err := dbConnection.GetWavs()
	if err != nil {
		return nil, fmt.Errorf("failed retrieving wav details from db: %w", err)
	}

	return result, nil
}

// func getData() (model.WavFilesDetailsSlice, error) {

// 	cfg := mysql.Config{
// 		User:                 "root", //"user",
// 		Passwd:               "password",
// 		Net:                  "tcp",
// 		Addr:                 "audio_db",
// 		DBName:               "audio_db",
// 		AllowNativePasswords: true,
// 	}

// 	db, err := sql.Open("mysql", cfg.FormatDSN())
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open database: %w", err)
// 	}
// 	defer db.Close()

// 	err = db.Ping()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to contact database: %w", err)
// 	}

// 	queryString := fmt.Sprintf("SELECT name, length_seconds FROM audio_db.wavs;")
// 	rows, err := db.Query(queryString)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to query database: %w", err)
// 	}

// 	var wavs model.WavFilesDetailsSlice
// 	for rows.Next() {
// 		var wav model.WavFileDetails
// 		if err := rows.Scan(&wav.Name, &wav.Duration); err != nil {
// 			return wavs, fmt.Errorf("error reading database: %w", err)
// 		}
// 		wavs = append(wavs, wav)
// 	}

// 	return wavs, nil

// }

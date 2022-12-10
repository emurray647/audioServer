package processing

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emurray647/audioServer/internal/model"
	"github.com/go-sql-driver/mysql"
)

func List(w http.ResponseWriter, r *http.Request) {
	// get all elements that match query

	// var wavs model.WavFiles
	wavs, err := getData()
	if err != nil {
		http.Error(w, "failed to read from database", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&wavs)
	if err != nil {
		http.Error(w, "failed to encode wav data", http.StatusInternalServerError)
		return
	}

}

func getData() (model.WavFilesDetailsSlice, error) {

	cfg := mysql.Config{
		User:                 "root", //"user",
		Passwd:               "password",
		Net:                  "tcp",
		Addr:                 "audio_db",
		DBName:               "audio_db",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to contact database: %w", err)
	}

	queryString := fmt.Sprintf("SELECT name, length_seconds FROM audio_db.wavs;")
	rows, err := db.Query(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}

	var wavs model.WavFilesDetailsSlice
	for rows.Next() {
		var wav model.WavFileDetails
		if err := rows.Scan(&wav.Name, &wav.Duration); err != nil {
			return wavs, fmt.Errorf("error reading database: %w", err)
		}
		wavs = append(wavs, wav)
	}

	return wavs, nil

}

package dbconnector

import (
	"database/sql"
	"fmt"

	"github.com/emurray647/audioServer/internal/model"
	"github.com/go-sql-driver/mysql"
)

type DBConnection struct {
	DB *sql.DB
}

func OpenDBConnection() (*DBConnection, error) {
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

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to contact database: %w", err)
	}

	connection := &DBConnection{
		DB: db,
	}
	return connection, nil
}

func (dc *DBConnection) Close() {
	dc.DB.Close()
}

func (dc *DBConnection) CountWavFiles(name string) (int, error) {
	queryString := fmt.Sprintf("SELECT count(*) FROM audio_db.wavs WHERE name='%s'", name)
	row := dc.DB.QueryRow(queryString)
	var count int
	err := row.Scan(&count)
	return count, err
}

func (dc *DBConnection) AddWavFile(wav *model.WavFile) error {
	_, err := dc.DB.Exec("INSERT INTO audio_db.wavs "+
		"(name, file_size, length_seconds, num_channels, sample_rate, audio_format, avg_bytes_per_sec, file_uri) "+
		"VALUES (?,?,?,?,?,?,?,?);",
		wav.Name, wav.FileSize, wav.Duration, wav.NumChannels, wav.SampleRate, wav.AudioFormat, wav.AvgBytesPerSec, wav.URI)
	if err != nil {
		return fmt.Errorf("failed to execute SQL statement: %w", err)
	}
	return nil
}

func (dc *DBConnection) DeleteWav(name string) (string, error) {
	tx, err := dc.DB.Begin()
	if err != nil {
		return "", fmt.Errorf("failed creating db transaction")
	}
	defer tx.Rollback()

	queryString := fmt.Sprintf("SELECT file_uri FROM audio_db.wavs WHERE name='%s';", name)
	row := tx.QueryRow(queryString)
	var fileURI string
	if err = row.Scan(&fileURI); err != nil {
		return "", err
	}

	execString := fmt.Sprintf("DELETE FROM audio_db.wavs WHERE name='%s'", name)
	_, err = tx.Exec(execString)
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction")
	}

	return fileURI, err
}

func (dc *DBConnection) GetWavURI(name string) (string, error) {
	queryString := fmt.Sprintf("SELECT file_uri FROM audio_db.wavs WHERE name='%s'", name)
	row := dc.DB.QueryRow(queryString)

	var uri string
	err := row.Scan(&uri)
	return uri, err
}

func (dc *DBConnection) GetWavDetails(name string) (*model.WavFileDetails, error) {

	queryString := fmt.Sprintf("SELECT "+
		"name, file_size, length_seconds, num_channels, sample_rate, audio_format, avg_bytes_per_sec "+
		"FROM audio_db.wavs WHERE name='%s'", name)
	row := dc.DB.QueryRow(queryString)

	var details model.WavFileDetails

	err := row.Scan(&details.Name, &details.FileSize, &details.Duration, &details.NumChannels, &details.SampleRate,
		&details.AudioFormat, &details.AvgBytesPerSec)

	if err != nil {
		return nil, fmt.Errorf("failed to read details from DB: %w", err)
	}

	return &details, nil
}

func (dc *DBConnection) GetWavs() (*model.WavFilesDetailsSlice, error) {
	queryString := fmt.Sprintf("SELECT " +
		"name, file_size, length_seconds, num_channels, sample_rate, audio_format, avg_bytes_per_sec" +
		" FROM audio_db.wavs")

	rows, err := dc.DB.Query(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed querying database: %w", err)
	}

	var result model.WavFilesDetailsSlice
	for rows.Next() {
		details := model.WavFileDetails{}
		err = rows.Scan(&details.Name, &details.FileSize, &details.Duration, &details.NumChannels,
			&details.SampleRate, &details.AudioFormat, &details.AvgBytesPerSec)
		if err != nil {
			return &result, fmt.Errorf("error scanning query result: %w", err)
		}
		result = append(result, details)
	}

	return &result, nil
}

package dbconnector

import (
	"database/sql"
	"fmt"

	"github.com/emurray647/audioServer/internal/model"
	"github.com/go-sql-driver/mysql"
)

const tableName = "audio_db.audio_files"

// A struct to handle all interactions with the database
type DBConnection struct {
	DB *sql.DB
}

// Opens a connection to the database
func OpenDBConnection(dbUser, dbPass, dbHost, dbName string) (*DBConnection, error) {
	cfg := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPass,
		Net:                  "tcp",
		Addr:                 dbHost,
		DBName:               dbName,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	connection := &DBConnection{
		DB: db,
	}
	return connection, nil
}

// Closes a database connection
func (dc *DBConnection) Close() {
	dc.DB.Close()
}

// Counts the number of entries in the table with the provided name
//    should be at most one
func (dc *DBConnection) CountFiles(name string) (int, error) {
	queryString := fmt.Sprintf("SELECT count(*) FROM %s WHERE name='%s'", tableName, name)
	row := dc.DB.QueryRow(queryString)
	var count int
	err := row.Scan(&count)
	return count, err
}

// Adds an audio file to the database
func (dc *DBConnection) AddFile(file *model.AudioFile) error {
	_, err := dc.DB.Exec("INSERT INTO "+tableName+
		" (name, file_size, format, duration, num_channels, sample_rate, avg_bytes_per_second, file_uri) "+
		"VALUES (?,?,?,?,?,?,?,?);",
		file.Name, file.FileSize, file.Format, file.Duration, file.NumChannels, file.SampleRate, file.AvgBytesPerSec, file.URI)
	if err != nil {
		return fmt.Errorf("failed to execute SQL statement: %w", err)
	}
	return nil
}

// Deletes an audio file entry from the database
func (dc *DBConnection) DeleteFile(name string) (string, error) {
	tx, err := dc.DB.Begin()
	if err != nil {
		return "", fmt.Errorf("failed creating db transaction")
	}
	defer tx.Rollback()

	queryString := fmt.Sprintf("SELECT file_uri FROM %s WHERE name='%s';", tableName, name)
	row := tx.QueryRow(queryString)
	var fileURI string
	if err = row.Scan(&fileURI); err != nil {
		return "", err
	}

	execString := fmt.Sprintf("DELETE FROM %s WHERE name='%s'", tableName, name)
	_, err = tx.Exec(execString)
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction")
	}

	return fileURI, err
}

// Gets the URI associated with the file name provided
func (dc *DBConnection) GetFileURI(name string) (string, error) {
	queryString := fmt.Sprintf("SELECT file_uri FROM %s WHERE name='%s'", tableName, name)
	row := dc.DB.QueryRow(queryString)

	var uri string
	err := row.Scan(&uri)
	return uri, err
}

// Retrieves the AudioFileDetails associated with the provided file name
func (dc *DBConnection) GetFileDetails(name string) (*model.AudioFileDetails, error) {
	queryString := fmt.Sprintf("SELECT "+
		"name, file_size, format, duration, num_channels, sample_rate, avg_bytes_per_second "+
		"FROM %s WHERE name='%s'", tableName, name)
	row := dc.DB.QueryRow(queryString)

	var details model.AudioFileDetails

	err := row.Scan(&details.Name, &details.FileSize, &details.Format, &details.Duration, &details.NumChannels,
		&details.SampleRate, &details.AvgBytesPerSec)

	if err != nil {
		return nil, fmt.Errorf("failed to read details from DB: %w", err)
	}

	return &details, nil
}

// Gets all the AudioFileDetails that satisfy the provided filters
func (dc *DBConnection) GetFiles(filterStrings []string) (*model.AudioFileDetailsSlice, error) {
	queryString := fmt.Sprintf("SELECT "+
		"name, file_size, format, duration, num_channels, sample_rate, avg_bytes_per_second"+
		" FROM %s", tableName)

	if len(filterStrings) > 0 {
		queryString += " WHERE " + filterStrings[0]
		for _, filter := range filterStrings[1:] {
			queryString += ", " + filter
		}
	}

	rows, err := dc.DB.Query(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed querying database: %w", err)
	}

	var result model.AudioFileDetailsSlice
	for rows.Next() {
		details := model.AudioFileDetails{}
		err = rows.Scan(&details.Name, &details.FileSize, &details.Format, &details.Duration, &details.NumChannels,
			&details.SampleRate, &details.AvgBytesPerSec)
		if err != nil {
			return &result, fmt.Errorf("error scanning query result: %w", err)
		}
		result = append(result, details)
	}

	return &result, nil
}

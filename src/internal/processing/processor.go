package processing

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
)

const (
	writePrefix = "/data"
)

func Upload(w http.ResponseWriter, r *http.Request) {

	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(buffer))

	name := r.URL.Query().Get("name")
	if name == "" {
		name = generateName(buffer)
	}

	filename := fmt.Sprintf("%s/%s.wav", writePrefix, name)
	err = os.WriteFile(filename, buffer, 0644)
	if err != nil {
		panic(err)
	}

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
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("INSERT INTO audio_db.wavs (name, length_seconds, file_url) VALUES (?,?,?)", name, 0, filename)
	if err != nil {
		panic(err)
	}
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandomName(size int) string {
	result := make([]byte, size)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func generateName(buffer []byte) string {
	hash := md5.Sum(buffer)
	return hex.EncodeToString(hash[:])
}

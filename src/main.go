package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/emurray647/audioServer/internal/dbconnector"
	"github.com/emurray647/audioServer/internal/routing"
)

var (
	dbUser = flag.String("dbuser", "user", "db username")
	dbPass = flag.String("dbpass", "password", "db password")
	dbHost = flag.String("dbhost", "audio_db", "db address")
	dbName = flag.String("dbname", "audio_db", "db name")

	fileDirectory = flag.String("fileDirectory", "/data", "directory to store uploaded files")
)

func main() {
	flag.Parse()

	// wait a couple seconds to give the DB a chance to come up
	time.Sleep(2 * time.Second)

	db, err := dbconnector.OpenDBConnection(*dbUser, *dbPass, *dbHost, *dbName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := routing.InitializeRoutes(db, *fileDirectory)

	http.ListenAndServe(":8080", router)
}

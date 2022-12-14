package main

import (
	"flag"
	"fmt"
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

	fmt.Printf("using user %s\n", *dbUser)

	time.Sleep(2 * time.Second)

	db, err := dbconnector.OpenDBConnection(*dbUser, *dbPass, *dbHost, *dbName)
	if err != nil {
		panic(err)
	}

	router := routing.InitializeRoutes(db, *fileDirectory)

	http.ListenAndServe(":8080", router)
}

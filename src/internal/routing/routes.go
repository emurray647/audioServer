package routing

import (
	"fmt"
	"net/http"

	"github.com/emurray647/audioServer/internal/processing"
	"github.com/gorilla/mux"
)

func InitializeRoutes() *mux.Router {
	r := mux.NewRouter()

	r.StrictSlash(false)

	r.HandleFunc("/files/", processing.Upload).Methods("POST")

	r.HandleFunc("/list/", processing.List).Methods("GET")

	r.HandleFunc("/download/", Downloads)

	return r
}

func Downloads(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got Downloads")
}

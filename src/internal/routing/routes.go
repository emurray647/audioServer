package routing

import (
	"fmt"
	"net/http"

	"github.com/emurray647/audioServer/internal/processing"
	"github.com/gorilla/mux"
)

func InitializeRoutes() *mux.Router {
	r := mux.NewRouter()

	r.StrictSlash(true)

	r.HandleFunc("/files/", processing.Upload).Methods("POST")
	r.HandleFunc("/files/{filename}", processing.Delete).Methods("DELETE")

	r.HandleFunc("/list/", processing.List).Methods("GET")

	r.HandleFunc("/download/", processing.Download).Methods("GET")

	return r
}

func Downloads(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got Downloads")
}

package routing

import (
	"net/http"
	"strings"

	"github.com/emurray647/audioServer/internal/processing"
	"github.com/gorilla/mux"
)

// func InitializeRoutes() *mux.Router {
func InitializeRoutes() http.Handler {
	r := mux.NewRouter()

	r.StrictSlash(true)

	r.HandleFunc("/files", processing.Upload).Methods("POST")
	r.HandleFunc("/files/{filename}", processing.Delete).Methods("DELETE")

	// r.HandleFunc("/list/", processing.List).Methods("GET")
	r.HandleFunc("/list", processing.CreateListHandler()).Methods("GET")

	r.HandleFunc("/download", processing.Download).Methods("GET")

	r.HandleFunc("/info", processing.Info).Methods("GET")

	return middleware(r)
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		next.ServeHTTP(w, r)
	})
}

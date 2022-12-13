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

	// POST /files
	r.HandleFunc("/files", processing.Upload).Methods("POST")

	// DELETE /files
	r.HandleFunc("/files", processing.Delete).Methods("DELETE")

	// GET /list
	r.HandleFunc("/list", processing.CreateListHandler()).Methods("GET")

	// GET /download
	r.HandleFunc("/download", processing.Download).Methods("GET")

	// GET /info
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

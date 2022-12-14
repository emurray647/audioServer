package routing

import (
	"net/http"
	"strings"

	"github.com/emurray647/audioServer/internal/dbconnector"
	"github.com/emurray647/audioServer/internal/processing"
	"github.com/gorilla/mux"
)

// func InitializeRoutes() *mux.Router {
func InitializeRoutes(db *dbconnector.DBConnection, filePrefix string) http.Handler {
	r := mux.NewRouter()

	processor := processing.NewRequestProcessor(db, filePrefix)

	// POST /files
	r.HandleFunc("/files", processor.Upload).Methods("POST")

	// DELETE /files
	r.HandleFunc("/files", processor.Delete).Methods("DELETE")

	// GET /list
	r.HandleFunc("/list", processor.CreateListHandler()).Methods("GET")

	// GET /download
	r.HandleFunc("/download", processor.Download).Methods("GET")

	// GET /info
	r.HandleFunc("/info", processor.Info).Methods("GET")

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

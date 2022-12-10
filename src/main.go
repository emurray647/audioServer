package main

import (
	"fmt"
	"net/http"

	"github.com/emurray647/audioServer/internal/routing"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Hello")

	// routes()

	router := routing.InitializeRoutes()

	http.ListenAndServe(":8080", router)
}

func routes() {
	r := mux.NewRouter()

	r.HandleFunc("/download/", Downloads)

	// http.Handle("/", r)

	http.ListenAndServe(":8080", r)

	// 	srv := &http.Server{
	// 		Handler:      r,
	// 		Addr:         ":8080",
	// 		WriteTimeout: 5 * time.Second,
	// 		ReadTimeout:  5 * time.Second,
	// 	}

	// 	srv.ListenAndServe()
}

func Downloads(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got Downloads")
}

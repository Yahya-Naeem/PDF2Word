package main

import (
	"backend/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Register routes
	r.HandleFunc("/upload", handlers.UploadFileHandler).Methods("POST")
	r.HandleFunc("/download/{fileName}", handlers.DownloadFileHandler).Methods("GET") // ✅ Fixed route

	// Start server
	fmt.Println("✅ Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

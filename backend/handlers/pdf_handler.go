package handlers

import (
	"backend/services"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// Enable CORS
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w) // Allow frontend requests

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		http.Error(w, "❌ File upload error", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "❌ Failed to retrieve file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	uploadDir := "storage"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		_ = os.Mkdir(uploadDir, os.ModePerm)
	}

	pdfPath := filepath.Join(uploadDir, handler.Filename)
	dst, err := os.Create(pdfPath)
	if err != nil {
		http.Error(w, "❌ Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "❌ Error saving file", http.StatusInternalServerError)
		return
	}

	services.ConvertPDFToWordWithPython(pdfPath)
	wordFilePath := pdfPath[:len(pdfPath)-4] + ".docx"

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "http://localhost:8080/download/%s", filepath.Base(wordFilePath))
}

// Handles file download requests
func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract filename from URL path
	vars := mux.Vars(r)
	fileName := vars["fileName"] // ✅ Correct way to get path parameter

	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	// Define file path (assuming files are stored in "storage" folder)
	filePath := filepath.Join("storage", fileName)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set headers for file download
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")

	// Copy file contents to response
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
	}
}

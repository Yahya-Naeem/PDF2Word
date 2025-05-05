package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"pdf2word/services"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm((10 << 20))
	if err != nil {
		http.Error(w, "File upload error", http.StatusBadRequest)
		return
	}

	file, Handler, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	uploadDir := "storage"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
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

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	fileName := filepath.Base(r.URL.Path)
	filePath := filepath.Join("storage", fileName)

	http.ServeFile(w, r, filePath)
}

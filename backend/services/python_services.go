package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Convert PDF to Word using Python script
func ConvertPDFToWordWithPython(pdfPath string) {
	// Ensure "storage" directory exists
	storageDir := "storage"
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		_ = os.Mkdir(storageDir, os.ModePerm)
	}

	// Generate output path
	outputPath := filepath.Join(storageDir, filepath.Base(pdfPath[:len(pdfPath)-4]+".docx"))

	// Define the Python script path
	scriptPath := "scripts/pdf_to_word.py"

	// Execute the Python script
	cmd := exec.Command("python", scriptPath, pdfPath, outputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("❌ Error running Python script: %v", err)
		return
	}

	fmt.Printf("✅ Converted with layout preserved: %s -> %s\n", pdfPath, outputPath)
}

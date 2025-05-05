package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"pdf2word/services"
	"strings" // Added for TrimSpace
)

func main() {
	fmt.Println("Enter the file path to be converted:")
	reader := bufio.NewReader(os.Stdin)
	filePath, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading file path: %v", err)
	}

	// Trim whitespace and clean file path
	filePath = strings.TrimSpace(filePath)
	filePath = filepath.Clean(filePath)

	fmt.Println("read pdf function firing forom here \n")
	//services.readPdf(filePath)
	// Get absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatalf("Error getting absolute file path: %v", err)
	}

	// Check if the file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatalf("File not found: %s", absPath)
	}

	fmt.Printf("File found at %s\n", absPath)

	// Convert PDF to Word
	services.ConvertPDFToWordWithPython(absPath)
}

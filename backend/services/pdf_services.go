package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ledongthuc/pdf"
	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
)

// Initialize UniOffice license from an environment variable
func init() {
	// Directly declare your API key here
	var apiKey string = "7069264cb1095366af05c8fb660e728d9641b8b5523a82c92941b7a20b74c242"
	err := license.SetMeteredKey(apiKey)
	if err != nil {
		log.Fatal("❌ Failed to set API key:", err)
	}
	fmt.Println("✅ API key loaded successfully")
}

// Extract formatted text from a PDF file
func ExtractFormattedTextFromPDF(filePath string) ([]string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var extractedText []string
	totalPages := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			var rowText string
			for _, word := range row.Content {
				rowText += word.S + " "
			}
			extractedText = append(extractedText, rowText)
		}
	}
	return extractedText, nil
}

// Save extracted text as a formatted Word document (.docx)
func saveFormattedTextToWord(textLines []string, outputPath string) error {
	doc := document.New()
	defer doc.Close()

	for _, line := range textLines {
		para := doc.AddParagraph()
		run := para.AddRun()
		run.AddText(line)

		// Example: Making the first line bold for demonstration
		if len(textLines) > 0 && line == textLines[0] {
			run.Properties().Bold()
		}
	}

	// Save to file
	err := doc.SaveToFile(outputPath)
	if err != nil {
		return err
	}
	return nil
}

// Convert formatted PDF to Word document
func ConvertFormattedPDFToWord(pdfPath string) {
	textLines, err := ExtractFormattedTextFromPDF(pdfPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Ensure "storage" directory exists
	storageDir := "storage"
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		_ = os.Mkdir(storageDir, os.ModePerm)
	}

	// Generate Word output path
	outputPath := filepath.Join(storageDir, filepath.Base(pdfPath[:len(pdfPath)-4]+".docx"))

	// Save extracted text to Word
	err = saveFormattedTextToWord(textLines, outputPath)
	if err != nil {
		log.Printf("❌ Failed to save Word file for %s: %v\n", pdfPath, err)
		return
	}
	fmt.Printf("✅ Converted with formatting: %s -> %s\n", pdfPath, outputPath)
}

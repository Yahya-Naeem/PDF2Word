package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
)

// Initialize UniOffice license from an environment variable
func init() {
	// Directly declare your API key here
	var apiKey string = ""
	err := license.SetMeteredKey(apiKey)
	if err != nil {
		log.Fatal("❌ Failed to set API key:", err)
	}
	fmt.Println("✅ API key loaded successfully")
}

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
			extractedText = append(extractedText, strings.TrimSpace(rowText))
		}
	}
	return extractedText, nil
}

func saveFormattedTextToWord(textLines []string, outputPath string) error {
	doc := document.New()
	defer doc.Close()

	// Add styles for different text types
	styles := document.NewStyleSheet()
	headingStyle := styles.AddParagraphStyle("Heading1")
	headingStyle.SetRunProperties(document.NewRunProperties())
	headingStyle.RunProperties().SetBold(true)
	headingStyle.RunProperties().SetSize(24)
	headingStyle.RunProperties().SetColor(document.Color{R: 0, G: 0, B: 0})

	bulletStyle := styles.AddParagraphStyle("Bullet")
	bulletStyle.SetRunProperties(document.NewRunProperties())
	bulletStyle.RunProperties().SetSize(12)
	bulletStyle.SetNumberingStyle(document.NewNumberingStyle())
	bulletStyle.NumberingStyle().SetType(document.ST_NumberFormatBullet)

	// Process each line
	for i, line := range textLines {
		para := doc.AddParagraph()

		// Detect and apply formatting
		if i == 0 {
			// First line is heading
			para.SetStyle("Heading1")
		} else if strings.HasPrefix(line, "-") {
			// Lines starting with - are bullets
			para.SetStyle("Bullet")
		} else {
			// Regular paragraph
			para.SetStyle("Normal")
		}

		run := para.AddRun()
		run.AddText(line)
	}

	return doc.SaveToFile(outputPath)
}

// Heuristic to detect headings
func isHeading(line string) bool {
	// Check if the line is in ALL CAPS or ends with a colon
	return strings.ToUpper(line) == line || strings.HasSuffix(line, ":")
}

// Heuristic to detect bullet points
func isBullet(line string) bool {
	// Check if the line starts with a bullet character (•, -, *, etc.)
	return strings.HasPrefix(line, "•") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*")
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

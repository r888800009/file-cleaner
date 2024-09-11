package file_cleaner

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/cheggaaa/go-poppler"
	"github.com/heussd/pdftotext-go"
)

func IsPDF(path string) (bool, error) {
	// read the first 4 bytes of the file
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	header := make([]byte, 4)
	_, err = file.Read(header)
	if err != nil {
		return false, err
	}

	// check if the header is a PDF
	if string(header) != "%PDF" {
		return false, nil
	}

	return true, nil
}

// https://github.com/heussd/pdftotext-go
func extractPdfContent(pdfPath string) (string, error) {
	var buf bytes.Buffer

	pdf, err := os.ReadFile(pdfPath)
	if err != nil {
		return "", err
	}

	pages, err := pdftotext.Extract(pdf)
	if err != nil {
		return "", err
	}

	// for each page, it returns a list of Page objects
	for _, page := range pages {
		buf.WriteString(page.Content)
	}

	return buf.String(), nil
}

// https://github.com/jacobkring/go-poppler
func extractPdfGOPoppler(pdfPath string) (string, error) {
	doc, err := poppler.Open(pdfPath)
	if err != nil {
		return "", err
	}
	defer doc.Close()

	numPages := doc.GetNPages()
	var buf bytes.Buffer
	for i := 0; i < numPages; i++ {
		page := doc.GetPage(i)
		text := page.Text()

		buf.WriteString(text)
	}
	return buf.String(), nil
}

func extractPdfCommand(path string) (string, error) {
	args := []string{path, "-"}
	output, err := exec.Command("pdftotext", args...).Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func ExtractPdf(path string) (string, error) {
	// we have been evaluating the performance of different libraries, and extractPdfGOPoppler is the fastest
	// so we will remove the other implementations you can find them in the git history
	// return extractPdfContent(pdfPath)
	// return extractPdfCommand(pdfPath)
	return extractPdfGOPoppler(path)
}

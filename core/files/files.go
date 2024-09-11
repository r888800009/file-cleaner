package file_cleaner

import (
	"bytes"
	"os"

	"github.com/cheggaaa/go-poppler"
)

func IsPDF(path string) (bool, error) {
	// read the first 4 bytes of the file
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// check if file < 4 bytes
	stat, err := file.Stat()
	if err != nil {
		return false, err
	}

	if stat.Size() < 4 {
		return false, nil
	}

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

func ExtractPdf(path string) (string, error) {
	// we have been evaluating the performance of different libraries, and extractPdfGOPoppler is the fastest
	// so we will remove the other implementations you can find them in the git history r888800009/file_cleaner@e200970
	return extractPdfGOPoppler(path)
}

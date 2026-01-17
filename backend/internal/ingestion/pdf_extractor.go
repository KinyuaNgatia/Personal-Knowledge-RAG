package ingestion

import (
	"github.com/ledongthuc/pdf"
)

// ExtractTextByPage opens a PDF file and extracts text from each page.
// It returns a map where the key is the page number and the value is the text.
func ExtractTextByPage(path string) (map[int]string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pages := make(map[int]string)
	for pageNum := 1; pageNum <= r.NumPage(); pageNum++ {
		p := r.Page(pageNum)
		if p.V.IsNull() {
			continue
		}

		text, err := p.GetPlainText(nil)
		if err != nil {
			return nil, err
		}
		pages[pageNum] = text
	}

	return pages, nil
}

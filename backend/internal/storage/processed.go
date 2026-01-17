package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ProcessedDocument struct {
	DocumentID string      `json:"document_id"`
	Source     string      `json:"source"`
	Pages      interface{} `json:"pages"`
}

func SaveProcessedDocument(docID string, source string, pages interface{}) error {
	path := filepath.Join("data", "processed", docID+".json")

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(ProcessedDocument{
		DocumentID: docID,
		Source:     source,
		Pages:      pages,
	})
}

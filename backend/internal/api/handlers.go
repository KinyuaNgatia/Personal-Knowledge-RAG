package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/KinyuaNgatia/Personal-Knowledge-RAG/backend/internal/ingestion"
	"github.com/KinyuaNgatia/Personal-Knowledge-RAG/backend/internal/storage"
	"github.com/google/uuid"
)

const maxUploadSize = 10 << 20 // 10 MB

func IngestPDFHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "file too large or invalid form data",
		})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "file is required",
		})
		return
	}
	defer file.Close()

	if header.Header.Get("Content-Type") != "application/pdf" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "invalid file type, only PDFs are allowed",
		})
		return
	}

	documentID := uuid.New().String()
	filename := documentID + ".pdf"

	// Ensure raw and processed directories exist
	if err := os.MkdirAll(filepath.Join("data", "raw"), 0755); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "failed to create raw storage directory",
		})
		return
	}
	if err := os.MkdirAll(filepath.Join("data", "processed"), 0755); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "failed to create processed storage directory",
		})
		return
	}

	// Save the raw PDF
	dstPath := filepath.Join("data", "raw", filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "failed to store file",
		})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "failed to write file",
		})
		return
	}

	// Extract PDF text page by page ---
	pages, err := ingestion.ExtractTextByPage(dstPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "failed to extract text from PDF",
		})
		return
	}

	if err := storage.SaveProcessedDocument(documentID, header.Filename, pages); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "failed to store extracted text",
		})
		return
	}

	// --- Respond ---
	resp := IngestPDFResponse{
		DocumentID:    documentID,
		Filename:      header.Filename,
		Status:        "stored",
		ChunksCreated: 0,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

	_ = time.Now() // placeholder
}

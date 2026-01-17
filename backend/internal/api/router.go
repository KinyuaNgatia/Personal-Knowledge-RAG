package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/ingest/pdf", IngestPDFHandler)
	mux.HandleFunc("/api/chat", ChatHandler)
}

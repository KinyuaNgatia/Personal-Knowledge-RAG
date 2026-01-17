package api

import (
	"encoding/json"
	"net/http"
)

func IngestPDFHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: "not implemented",
	})
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: "not implemented",
	})
}

package api

type IngestPDFResponse struct {
	DocumentID    string `json:"document_id"`
	Filename      string `json:"filename"`
	Status        string `json:"status"`
	ChunksCreated int    `json:"chunks_created"`
}

type ChatRequest struct {
	Query string `json:"query"`
}

type Source struct {
	Document string `json:"document"`
	Page     int    `json:"page"`
	Snippet  string `json:"snippet"`
}

type ChatResponse struct {
	Answer  string   `json:"answer"`
	Sources []Source `json:"sources"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

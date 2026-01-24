package handler

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/luckydiss/ocr_tgapp/internal/gemini"
)

const maxImageSize = 10 * 1024 * 1024 // 10MB

type ExtractRequest struct {
	ImageBase64 string `json:"image_base64"`
	MimeType    string `json:"mime_type"`
}

type ExtractResponse struct {
	Success bool            `json:"success"`
	Payload *ExtractPayload `json:"payload"`
	Error   *string         `json:"error"`
}

type ExtractPayload struct {
	Markdown string `json:"markdown"`
}

type ExtractHandler struct {
	geminiService gemini.GeminiService
}

func NewExtractHandler(gs gemini.GeminiService) *ExtractHandler {
	return &ExtractHandler{geminiService: gs}
}

func (h *ExtractHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ExtractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.ImageBase64 == "" {
		sendError(w, "image_base64 is required", http.StatusBadRequest)
		return
	}

	if req.MimeType == "" {
		sendError(w, "mime_type is required", http.StatusBadRequest)
		return
	}

	imgData, err := base64.StdEncoding.DecodeString(req.ImageBase64)
	if err != nil {
		sendError(w, "invalid base64 encoding", http.StatusBadRequest)
		return
	}

	if len(imgData) > maxImageSize {
		sendError(w, "image too large (max 10MB)", http.StatusBadRequest)
		return
	}

	markdown, err := h.geminiService.ExtractMarkdown(r.Context(), imgData, req.MimeType)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccess(w, markdown)
}

func sendError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ExtractResponse{
		Success: false,
		Payload: nil,
		Error:   &msg,
	})
}

func sendSuccess(w http.ResponseWriter, markdown string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ExtractResponse{
		Success: true,
		Payload: &ExtractPayload{Markdown: markdown},
		Error:   nil,
	})
}

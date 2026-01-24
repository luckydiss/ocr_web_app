package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luckydiss/ocr_tgapp/internal/gemini"
)

func TestExtractHandler_MissingImage(t *testing.T) {
	mock := gemini.NewMockGeminiService("", nil)
	handler := NewExtractHandler(mock)

	body := `{"mime_type": "image/png"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/extract", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}

	var resp ExtractResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Success {
		t.Error("expected success to be false")
	}
	if resp.Error == nil || *resp.Error != "image_base64 is required" {
		t.Errorf("unexpected error message: %v", resp.Error)
	}
}

func TestExtractHandler_MissingMimeType(t *testing.T) {
	mock := gemini.NewMockGeminiService("", nil)
	handler := NewExtractHandler(mock)

	imgBase64 := base64.StdEncoding.EncodeToString([]byte("fake-image"))
	body := `{"image_base64": "` + imgBase64 + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/extract", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestExtractHandler_InvalidBase64(t *testing.T) {
	mock := gemini.NewMockGeminiService("", nil)
	handler := NewExtractHandler(mock)

	body := `{"image_base64": "not-valid-base64!!!", "mime_type": "image/png"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/extract", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestExtractHandler_UnsupportedMimeType(t *testing.T) {
	mock := gemini.NewMockGeminiService("", nil)
	handler := NewExtractHandler(mock)

	imgBase64 := base64.StdEncoding.EncodeToString([]byte("fake-image"))
	body := `{"image_base64": "` + imgBase64 + `", "mime_type": "text/plain"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/extract", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}
}

func TestExtractHandler_Success(t *testing.T) {
	expectedMarkdown := "# Title\n\n$ E = mc^2 $"
	mock := gemini.NewMockGeminiService(expectedMarkdown, nil)
	handler := NewExtractHandler(mock)

	imgBase64 := base64.StdEncoding.EncodeToString([]byte("fake-image-data"))
	body := `{"image_base64": "` + imgBase64 + `", "mime_type": "image/png"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/extract", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var resp ExtractResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if !resp.Success {
		t.Error("expected success to be true")
	}
	if resp.Payload == nil || resp.Payload.Markdown != expectedMarkdown {
		t.Errorf("expected markdown %q, got %v", expectedMarkdown, resp.Payload)
	}
}

func TestExtractHandler_ServiceError(t *testing.T) {
	mock := gemini.NewMockGeminiService("", errors.New("API error"))
	handler := NewExtractHandler(mock)

	imgBase64 := base64.StdEncoding.EncodeToString([]byte("fake-image-data"))
	body := `{"image_base64": "` + imgBase64 + `", "mime_type": "image/png"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/extract", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}
}

func TestExtractHandler_MethodNotAllowed(t *testing.T) {
	mock := gemini.NewMockGeminiService("", nil)
	handler := NewExtractHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/extract", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rec.Code)
	}
}

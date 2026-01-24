package gemini

import (
	"context"
	"errors"
	"testing"
)

func TestMockGeminiService_ExtractMarkdown_Success(t *testing.T) {
	expectedResponse := "# Title\n\n$ E = mc^2 $"
	mock := NewMockGeminiService(expectedResponse, nil)

	result, err := mock.ExtractMarkdown(context.Background(), []byte("fake-image-data"), "image/png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expectedResponse {
		t.Errorf("expected %q, got %q", expectedResponse, result)
	}
}

func TestMockGeminiService_ExtractMarkdown_EmptyImage(t *testing.T) {
	mock := NewMockGeminiService("response", nil)

	_, err := mock.ExtractMarkdown(context.Background(), []byte{}, "image/png")
	if err == nil {
		t.Fatal("expected error for empty image")
	}

	if err.Error() != "image data is empty" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestMockGeminiService_ExtractMarkdown_InvalidMimeType(t *testing.T) {
	mock := NewMockGeminiService("response", nil)

	_, err := mock.ExtractMarkdown(context.Background(), []byte("data"), "text/plain")
	if err == nil {
		t.Fatal("expected error for invalid mime type")
	}
}

func TestMockGeminiService_ExtractMarkdown_ServiceError(t *testing.T) {
	expectedErr := errors.New("API error")
	mock := NewMockGeminiService("", expectedErr)

	_, err := mock.ExtractMarkdown(context.Background(), []byte("data"), "image/png")
	if err == nil {
		t.Fatal("expected error from service")
	}

	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}

func TestIsValidMimeType(t *testing.T) {
	tests := []struct {
		mimeType string
		valid    bool
	}{
		{"image/jpeg", true},
		{"image/png", true},
		{"image/gif", true},
		{"image/webp", true},
		{"text/plain", false},
		{"application/pdf", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.mimeType, func(t *testing.T) {
			if got := isValidMimeType(tt.mimeType); got != tt.valid {
				t.Errorf("isValidMimeType(%q) = %v, want %v", tt.mimeType, got, tt.valid)
			}
		})
	}
}

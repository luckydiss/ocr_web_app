package gemini

import (
	"context"
	"fmt"
)

type MockGeminiService struct {
	Response string
	Error    error
}

func NewMockGeminiService(response string, err error) *MockGeminiService {
	return &MockGeminiService{
		Response: response,
		Error:    err,
	}
}

func (m *MockGeminiService) ExtractMarkdown(ctx context.Context, imgData []byte, mimeType string) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}

	if len(imgData) == 0 {
		return "", fmt.Errorf("image data is empty")
	}

	if !isValidMimeType(mimeType) {
		return "", fmt.Errorf("unsupported mime type: %s", mimeType)
	}

	return m.Response, nil
}

package gemini

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GeminiService interface {
	ExtractMarkdown(ctx context.Context, imgData []byte, mimeType string) (string, error)
}

type geminiService struct {
	apiKey   string
	endpoint string
	model    string
	client   *http.Client
}

const extractionPrompt = `Extract all text and mathematical formulas from this image.
Return the result in Markdown format.
Use LaTeX syntax wrapped in $ for inline math and $$ for display math.
Preserve document structure (headings, lists, paragraphs).
Only return the extracted content, no explanations or preamble.`

func NewGeminiService(ctx context.Context, apiKey, endpoint string) (GeminiService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}

	if endpoint == "" {
		endpoint = "http://127.0.0.1:8045"
	}

	return &geminiService{
		apiKey:   apiKey,
		endpoint: endpoint,
		model:    "gemini-2.5-flash",
		client:   &http.Client{},
	}, nil
}

type geminiRequest struct {
	Contents []content `json:"contents"`
}

type content struct {
	Role  string `json:"role"`
	Parts []part `json:"parts"`
}

type part struct {
	Text       string      `json:"text,omitempty"`
	InlineData *inlineData `json:"inline_data,omitempty"`
}

type inlineData struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
}

type geminiResponse struct {
	Candidates []candidate `json:"candidates"`
	Error      *apiError   `json:"error,omitempty"`
}

type candidate struct {
	Content contentResponse `json:"content"`
}

type contentResponse struct {
	Parts []partResponse `json:"parts"`
}

type partResponse struct {
	Text string `json:"text"`
}

type apiError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (s *geminiService) ExtractMarkdown(ctx context.Context, imgData []byte, mimeType string) (string, error) {
	if len(imgData) == 0 {
		return "", fmt.Errorf("image data is empty")
	}

	if !isValidMimeType(mimeType) {
		return "", fmt.Errorf("unsupported mime type: %s", mimeType)
	}

	reqBody := geminiRequest{
		Contents: []content{
			{
				Role: "user",
				Parts: []part{
					{
						InlineData: &inlineData{
							MimeType: mimeType,
							Data:     base64.StdEncoding.EncodeToString(imgData),
						},
					},
					{
						Text: extractionPrompt,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", s.endpoint, s.model, s.apiKey)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if geminiResp.Error != nil {
		return "", fmt.Errorf("API error: %s", geminiResp.Error.Message)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	var result string
	for _, p := range geminiResp.Candidates[0].Content.Parts {
		result += p.Text
	}

	return result, nil
}

func isValidMimeType(mimeType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	return validTypes[mimeType]
}

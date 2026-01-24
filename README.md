# OCR to Markdown - Telegram Mini App

A Telegram Mini App that converts images containing text and mathematical formulas to Markdown using Google Gemini API.

## Features

- Extract text from images (JPEG, PNG, GIF, WebP)
- Recognize mathematical formulas and convert to LaTeX
- Render Markdown with math support (KaTeX)
- Copy extracted Markdown to clipboard
- Native Telegram Mini App integration

## Project Structure

```
ocr_tgapp/
├── cmd/server/main.go          # Application entrypoint
├── internal/
│   ├── config/config.go        # Configuration loading
│   ├── gemini/
│   │   ├── service.go          # Gemini API service
│   │   ├── mock.go             # Mock for testing
│   │   └── service_test.go     # Unit tests
│   └── handler/
│       ├── extract.go          # HTTP handler
│       ├── extract_test.go     # Handler tests
│       └── middleware.go       # CORS, logging
├── web/
│   ├── index.html              # TMA main page
│   ├── css/style.css           # Styles
│   └── js/
│       ├── app.js              # Main app logic
│       ├── telegram.js         # Telegram WebApp SDK
│       └── markdown.js         # Markdown/LaTeX rendering
├── Dockerfile
├── .env.example
└── README.md
```

## Prerequisites

- Go 1.22+
- Google Gemini API key

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/luckydiss/ocr_tgapp.git
   cd ocr_tgapp
   ```

2. Copy environment file and configure:
   ```bash
   cp .env.example .env
   # Edit .env and add your GEMINI_API_KEY
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

## Running Locally

```bash
# Set environment variables
export GEMINI_API_KEY=your_api_key

# Run the server
go run ./cmd/server

# Server starts at http://localhost:8080
```

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...
```

## API Endpoint

### POST /api/v1/extract

Extract text and formulas from an image.

**Request:**
```json
{
  "image_base64": "<base64_encoded_image>",
  "mime_type": "image/png"
}
```

**Response (Success):**
```json
{
  "success": true,
  "payload": {
    "markdown": "# Title\n\n$ E = mc^2 $"
  },
  "error": null
}
```

**Response (Error):**
```json
{
  "success": false,
  "payload": null,
  "error": "error message"
}
```

## Docker Deployment

```bash
# Build the image
docker build -t ocr-tgapp .

# Run the container
docker run -d \
  -p 8080:8080 \
  -e GEMINI_API_KEY=your_api_key \
  --name ocr-tgapp \
  ocr-tgapp
```

## Telegram Bot Setup

1. Create a bot with [@BotFather](https://t.me/BotFather)
2. Set the Menu Button URL to your deployed app URL
3. Configure the bot's webapp using:
   ```
   /setmenubutton
   ```

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| GEMINI_API_KEY | Yes | - | Google Gemini API key |
| PORT | No | 8080 | Server port |
| ALLOWED_ORIGINS | No | * | CORS allowed origins |

## License

MIT

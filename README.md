# SHODAN Backend

A Go backend service that acts as a gateway to the Python FastAPI model service. Built with Gin framework and GORM for SQLite storage.

## Architecture

```
Client -> Go Backend (Gin) -> Python FastAPI (Models)
                 ↓
           SQLite Database (Results)
```

## Setup

### Prerequisites
- Go 1.22 or later
- Python FastAPI service running on `http://localhost:8000`

### Installation

1. Navigate to the backend directory:
```bash
cd shodan-backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the application:
```bash
go build
```

4. Run the server:
```bash
./shodan-backend
# or
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### 1. Single Text Analysis
**POST** `/analyze/text`

Analyzes sentiment for a single text string.

```bash
curl -X POST http://localhost:8080/analyze/text \
  -H "Content-Type: application/json" \
  -d '{"text": "I love this product!"}'
```

### 2. Batch Text Analysis
**POST** `/analyze/batch`

Analyzes sentiment for multiple text strings.

```bash
curl -X POST http://localhost:8080/analyze/batch \
  -H "Content-Type: application/json" \
  -d '{"texts": ["I love this product!", "This is terrible"]}'
```

### 3. Text Summarization
**POST** `/summarize`

Summarizes long text content.

```bash
curl -X POST http://localhost:8080/summarize \
  -H "Content-Type: application/json" \
  -d '{"text": "Your long text content here..."}'
```

### 4. Results History
**GET** `/history`

Retrieves all past analysis results from the database.

```bash
curl http://localhost:8080/history
```

## Response Format

All endpoints return JSON responses:

### Success Response
```json
{
  "result": {
    "sentiment": "positive",
    "confidence": 0.95,
    "model": "bert-base-uncased"
  }
}
```

### Error Response
```json
{
  "error": "error message",
  "details": "detailed error information"
}
```

## Database

The application uses SQLite (`results.db`) to store all analysis results with the following schema:

```sql
CREATE TABLE results (
    id INTEGER PRIMARY KEY,
    text TEXT,
    output TEXT,
    confidence REAL,
    model TEXT,
    created_at DATETIME
);
```

## Project Structure

```
shodan-backend/
├── main.go              # Application entry point
├── db.go                # Database initialization
├── go.mod               # Go module dependencies
├── models/
│   └── result.go        # GORM model for results
└── routes/
    ├── analyze.go       # Analysis route handlers
    └── summarize.go     # Summarization route handler
```

## Dependencies

- **Gin**: HTTP web framework
- **GORM**: ORM for database operations
- **SQLite**: Database storage

## Development

To run in development mode with auto-reload, you can use:

```bash
# Install air for live reload (optional)
go install github.com/cosmtrek/air@latest

# Run with live reload
air
```

## Error Handling

The application includes proper error handling for:
- Invalid JSON payloads
- FastAPI service unavailable
- Database connection issues
- Malformed responses from FastAPI

All errors are returned as JSON with appropriate HTTP status codes.
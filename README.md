# RISA - Reverse Image Search Agreggator

A RESTful API for reverse image search across multiple engines, optimized for finding people on social media and other websites. Developed with Go, the Echo framework, and PostgreSQL.

## Features

- **Multiple Search Engines**: Support for Yandex (synchronous) and FaceCrawler (asynchronous)
- **Image Upload**: Upload images for reverse search via multipart form data
- **Async Job Tracking**: Check status of asynchronous search operations
- **Security**: Built-in CORS, rate limiting, and security headers
- **Database**: PostgreSQL integration with GORM

## API Endpoints

### POST `/v1/image/upload`

Upload an image for reverse search.

**Request:**
- Method: `POST`
- Content-Type: `multipart/form-data`
- Parameters:
  - `file` (required): Image file
  - `engine` (required): Search engine name (`yandex` or `facecrawler`)

**Response:**

For synchronous engines (Yandex):
```json
{
  "engine": "yandex",
  "result": {...}
}
```

For asynchronous engines (FaceCrawler):
```json
{
  "engine": "facecrawler",
  "job_id": "unique-job-id"
}
```

**Status Codes:**
- `200 OK`: Success
- `400 Bad Request`: Invalid request or engine
- `500 Internal Server Error`: Server error

### GET `/v1/image/status`

Check the status of an asynchronous search job.

**Request:**
- Method: `GET`
- Query Parameters:
  - `engine` (required): Search engine name
  - `job_id` (required): Job ID returned from upload endpoint

**Response:**
```json
{
  "engine": "facecrawler",
  "result": {...}
}
```

**Status Codes:**
- `200 OK`: Success
- `400 Bad Request`: Missing parameters or invalid engine
- `500 Internal Server Error`: Server error

## Setup

### Prerequisites

- Go 1.24.3 or higher
- PostgreSQL database
- Environment variables configured

### Environment Variables

Create a `.env` file in the root directory:

```env
API_PORT=:8080
# Add your database and other configuration variables
```

### Installation

```bash
go mod download
go run cmd/api/server.go
```

## Supported Search Engines

- **Yandex**: Synchronous reverse image search
- **FaceCrawler**: Asynchronous reverse image search for face detection

## Security

The API includes:
- Rate limiting (5 requests per second)
- CORS support
- Security headers
- Request size limits (15s timeout)
- Gzip compression

## License

This project is for educational purposes.


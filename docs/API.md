# API Documentation

## Overview

This API follows RESTful principles and uses JSON for request and response payloads. All endpoints are prefixed with `/api/v1`.

## Authentication

Most endpoints require authentication using JWT tokens. Include the token in the Authorization header:

```
Authorization: Bearer <your-token>
```

## Common Response Formats

### Success Response
```json
{
    "data": {
        // Response data
    },
    "meta": {
        "page": 1,
        "per_page": 10,
        "total": 100
    }
}
```

### Error Response
```json
{
    "error": {
        "code": "INVALID_INPUT",
        "message": "Invalid input provided",
        "details": {
            "email": "must be a valid email address"
        }
    }
}
```

## Endpoints

### User Management

#### Create User
```http
POST /api/v1/users
```

Request:
```json
{
    "email": "user@example.com",
    "username": "johndoe",
    "password": "securepassword123"
}
```

Response:
```json
{
    "data": {
        "id": 1,
        "email": "user@example.com",
        "username": "johndoe",
        "created_at": "2023-12-19T10:00:00Z",
        "updated_at": "2023-12-19T10:00:00Z"
    }
}
```

#### Get User
```http
GET /api/v1/users/:id
```

Response:
```json
{
    "data": {
        "id": 1,
        "email": "user@example.com",
        "username": "johndoe",
        "created_at": "2023-12-19T10:00:00Z",
        "updated_at": "2023-12-19T10:00:00Z"
    }
}
```

#### Update User
```http
PUT /api/v1/users/:id
```

Request:
```json
{
    "username": "newusername",
    "email": "newemail@example.com"
}
```

Response:
```json
{
    "data": {
        "id": 1,
        "email": "newemail@example.com",
        "username": "newusername",
        "created_at": "2023-12-19T10:00:00Z",
        "updated_at": "2023-12-19T10:30:00Z"
    }
}
```

#### Delete User
```http
DELETE /api/v1/users/:id
```

Response:
```json
{
    "data": {
        "message": "User successfully deleted"
    }
}
```

### Authentication

#### Login
```http
POST /api/v1/auth/login
```

Request:
```json
{
    "email": "user@example.com",
    "password": "securepassword123"
}
```

Response:
```json
{
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIs...",
        "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
        "expires_in": 3600
    }
}
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
```

Request:
```json
{
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

Response:
```json
{
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIs...",
        "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
        "expires_in": 3600
    }
}
```

### Health Check

#### Get Health Status
```http
GET /health
```

Response:
```json
{
    "status": "healthy",
    "version": "1.0.0",
    "timestamp": "2023-12-19T10:00:00Z",
    "services": {
        "database": "up",
        "redis": "up"
    }
}
```

## Error Codes

| Code | Description |
|------|-------------|
| INVALID_INPUT | The request payload contains invalid data |
| UNAUTHORIZED | Authentication is required or failed |
| FORBIDDEN | The authenticated user lacks permission |
| NOT_FOUND | The requested resource was not found |
| CONFLICT | The request conflicts with existing data |
| INTERNAL_ERROR | An unexpected error occurred |

## Pagination

For endpoints that return lists, use the following query parameters:

- `page`: Page number (default: 1)
- `per_page`: Items per page (default: 10, max: 100)
- `sort`: Sort field (e.g., `created_at`)
- `order`: Sort order (`asc` or `desc`)

Example:
```http
GET /api/v1/users?page=2&per_page=20&sort=created_at&order=desc
```

## Filtering

Use query parameters for filtering:

```http
GET /api/v1/users?username=john&created_after=2023-01-01
```

Available filters:
- `username`: Filter by username
- `email`: Filter by email
- `created_after`: Filter by creation date (after)
- `created_before`: Filter by creation date (before)
- `status`: Filter by user status

## Rate Limiting

The API implements rate limiting with the following defaults:
- 100 requests per minute for authenticated users
- 20 requests per minute for unauthenticated users

Rate limit headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1702980000
```

## Versioning

The API uses URL versioning (e.g., `/api/v1/`). Breaking changes will result in a new version number.

## CORS

The API supports CORS for browser-based applications. Configure allowed origins in the configuration file.

## Monitoring

The following endpoints are available for monitoring:

- `/metrics`: Prometheus metrics
- `/health`: Health check status
- `/debug/pprof`: Performance profiling (protected) 
# Authentication Service Documentation

## Overview

This project implements a professional JWT-based authentication service using Go, Gin, and PostgreSQL (Supabase). The implementation follows Go best practices with proper error handling, security measures, and clean architecture.

## Features

- **User Registration**: Secure user registration with email validation and password hashing
- **User Login**: Authentication with JWT token generation
- **JWT Token Management**: Token validation, refresh, and expiration handling
- **Password Security**: bcrypt hashing for secure password storage
- **Middleware Protection**: Route protection with JWT middleware
- **Database Integration**: PostgreSQL connection with automatic migrations
- **CORS Support**: Cross-origin request handling for frontend integration

## Architecture

### Directory Structure
```
internal/
├── config/          # Configuration management
├── controllers/     # HTTP handlers
├── database/        # Database connection and migrations
├── middleware/      # Authentication middleware
├── models/          # Data models and request/response structures
├── router/          # Route definitions
└── services/
    ├── auth/        # Authentication business logic
    └── jwt/         # JWT token management
```

## API Endpoints

### Public Endpoints

#### POST /api/auth/register
Register a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "created_at": "2023-...",
    "updated_at": "2023-..."
  },
  "token": "jwt.token.here"
}
```

#### POST /api/auth/login
Authenticate user and receive JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "message": "Login successful",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "created_at": "2023-...",
    "updated_at": "2023-..."
  },
  "token": "jwt.token.here"
}
```

#### POST /api/auth/refresh
Refresh an existing JWT token.

**Request Body:**
```json
{
  "token": "existing.jwt.token"
}
```

**Response:**
```json
{
  "message": "Token refreshed successfully",
  "token": "new.jwt.token"
}
```

### Protected Endpoints

#### GET /api/profile
Get current user profile (requires authentication).

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "created_at": "2023-...",
    "updated_at": "2023-..."
  }
}
```

## Setup Instructions

### 1. Environment Configuration

Create a `.env` file based on `env.example`:

```bash
# Server Configuration
PORT=8080

# Database Configuration
DATABASE_URL=postgresql://postgres.qzuqwvvurdouhcvaybju:[YOUR-PASSWORD]@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Run the Server

```bash
go run ./cmd/server
```

The server will:
- Connect to the PostgreSQL database
- Run automatic migrations to create the users table
- Start listening on the configured port (default: 8080)

## Security Features

### Password Security
- **bcrypt Hashing**: All passwords are hashed using bcrypt with default cost
- **No Password Exposure**: Passwords are never returned in API responses
- **Minimum Length**: 6-character minimum password requirement

### JWT Security
- **HMAC-SHA256 Signing**: Tokens signed with secure algorithm
- **Expiration**: 24-hour token expiration (configurable)
- **Claims Validation**: Proper token validation with claims verification
- **Secure Headers**: Authorization header with Bearer token format

### Database Security
- **SQL Injection Protection**: Parameterized queries prevent SQL injection
- **Connection Pooling**: Efficient database connection management
- **UUID Primary Keys**: UUIDs used for user IDs for better security

## Error Handling

The API provides comprehensive error handling with appropriate HTTP status codes:

- `400 Bad Request`: Invalid request body or validation errors
- `401 Unauthorized`: Invalid credentials or expired tokens
- `404 Not Found`: User or resource not found
- `409 Conflict`: User already exists during registration
- `500 Internal Server Error`: Server-side errors

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
```

## Testing with cURL

### Register a new user:
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Login:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Access protected route:
```bash
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

## Production Considerations

1. **JWT Secret**: Use a strong, randomly generated secret key
2. **HTTPS**: Always use HTTPS in production
3. **Rate Limiting**: Implement rate limiting for auth endpoints
4. **Logging**: Add comprehensive logging for security events
5. **Token Blacklisting**: Consider implementing token blacklisting for logout
6. **Password Policies**: Implement stronger password requirements
7. **Email Verification**: Add email verification for registration
8. **Two-Factor Authentication**: Consider adding 2FA for enhanced security

## Dependencies

- `github.com/gin-gonic/gin`: HTTP web framework
- `github.com/golang-jwt/jwt/v5`: JWT implementation
- `github.com/lib/pq`: PostgreSQL driver
- `golang.org/x/crypto/bcrypt`: Password hashing
- `github.com/joho/godotenv`: Environment variable loading

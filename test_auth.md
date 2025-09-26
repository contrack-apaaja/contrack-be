# Testing the Authentication Service

## Setup Instructions

1. **Create your .env file** (copy from env.example):
```bash
cp env.example .env
```

2. **Update the DATABASE_URL** in .env with your actual Supabase password:
```
DATABASE_URL=postgresql://postgres.qzuqwvvurdouhcvaybju:YOUR_ACTUAL_PASSWORD@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres
```

3. **Start the server**:
```bash
go run ./cmd/server
```

## Test Commands

### 1. Health Check
```bash
curl http://localhost:8080/api/hello
```

### 2. Register a new user
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. Login with the user
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 4. Access protected profile endpoint
```bash
# Replace YOUR_JWT_TOKEN with the token from login response
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 5. Refresh token
```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "token": "YOUR_JWT_TOKEN"
  }'
```

## Expected Responses

All endpoints return JSON responses with appropriate HTTP status codes:
- 200/201 for success
- 400 for bad requests
- 401 for unauthorized
- 409 for conflicts (user already exists)
- 500 for server errors

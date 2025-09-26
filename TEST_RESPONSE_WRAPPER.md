# 🧪 Testing the Response Wrapper

## Quick Test Commands

### 1. Health Check (New Format)
```bash
curl http://localhost:8080/api/hello
```
**New Response:**
```json
{
  "status": "success",
  "message": "Hello from Contrack API!",
  "data": {
    "supabase_configured": true,
    "version": "1.0.0",
    "service": "authentication"
  }
}
```

### 2. Register User (New Format)
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@contrack.com",
    "password": "hackathon2024"
  }'
```
**New Response:**
```json
{
  "status": "success",
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "test@contrack.com",
      "created_at": "2024-09-26T10:30:00Z",
      "updated_at": "2024-09-26T10:30:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 3. Login (New Format)
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@contrack.com",
    "password": "hackathon2024"
  }'
```
**New Response:**
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "test@contrack.com",
      "created_at": "2024-09-26T10:30:00Z",
      "updated_at": "2024-09-26T10:30:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 4. Get Profile (New Format)
```bash
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```
**New Response:**
```json
{
  "status": "success",
  "message": "User profile retrieved successfully",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "test@contrack.com",
      "created_at": "2024-09-26T10:30:00Z",
      "updated_at": "2024-09-26T10:30:00Z"
    }
  }
}
```

## Error Response Examples

### Validation Error
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "invalid-email"}'
```
**Response:**
```json
{
  "status": "error",
  "message": "Validation failed",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": "Key: 'UserRegistrationRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
  }
}
```

### Unauthorized Error
```bash
curl -X GET http://localhost:8080/api/profile
```
**Response:**
```json
{
  "status": "error",
  "message": "Authorization header is required",
  "error": {
    "code": "UNAUTHORIZED",
    "details": null
  }
}
```

### Conflict Error
```bash
# Try to register same user twice
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@contrack.com",
    "password": "hackathon2024"
  }'
```
**Response:**
```json
{
  "status": "error",
  "message": "User already exists",
  "error": {
    "code": "CONFLICT",
    "details": null
  }
}
```

### Wrong Credentials
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@contrack.com",
    "password": "wrongpassword"
  }'
```
**Response:**
```json
{
  "status": "error",
  "message": "Invalid email or password",
  "error": {
    "code": "UNAUTHORIZED",
    "details": null
  }
}
```

## Key Benefits

### 1. Consistent Structure
All responses now follow the same format, making frontend integration easier.

### 2. Clear Status Indication
- `success`: Everything worked (2xx)
- `error`: Server/client error (4xx/5xx)
- `fail`: Request failed due to client input

### 3. Structured Error Information
Errors include:
- Human-readable message
- Error code for programmatic handling
- Optional details for debugging

### 4. Better Frontend Integration
Frontend developers can now handle responses consistently:

```typescript
interface APIResponse<T = any> {
  status: 'success' | 'error' | 'fail';
  message: string;
  data?: T;
  error?: {
    code: string;
    details?: any;
  };
}

// Usage
if (response.status === 'success') {
  // Handle success
  console.log(response.message);
  const data = response.data;
} else {
  // Handle error
  console.error(response.message);
  if (response.error?.code === 'VALIDATION_ERROR') {
    // Handle validation specifically
  }
}
```

## Testing Checklist

- ✅ All success responses have `status: "success"`
- ✅ All error responses have `status: "error"`
- ✅ All responses have meaningful messages
- ✅ Success responses include relevant data
- ✅ Error responses include error codes
- ✅ HTTP status codes match response status
- ✅ Validation errors provide helpful details

Your API now provides a consistent, professional response format! 🎉

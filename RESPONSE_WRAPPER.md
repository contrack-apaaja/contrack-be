# 📦 Response Wrapper Documentation

## Overview

The response wrapper provides a standardized format for all API responses in the Contrack backend. This ensures consistency across all endpoints and makes it easier for frontend developers to handle responses.

## Response Structure

All API responses follow this standardized format:

```json
{
  "status": "success|error|fail",
  "message": "Human readable message",
  "data": {}, // Optional: Present on success responses
  "error": {  // Optional: Present on error responses
    "code": "ERROR_CODE",
    "details": "Additional error information"
  }
}
```

### Status Values

- **`success`**: Request completed successfully (2xx status codes)
- **`error`**: Server error occurred (5xx status codes)
- **`fail`**: Client error - bad request, validation, etc. (4xx status codes)

## Available Response Functions

### Success Responses

#### `OKResponse(c, message, data)` - 200 OK
```go
utils.OKResponse(c, "User retrieved successfully", gin.H{
    "user": user,
})
```
**Response:**
```json
{
  "status": "success",
  "message": "User retrieved successfully",
  "data": {
    "user": {...}
  }
}
```

#### `CreatedResponse(c, message, data)` - 201 Created
```go
utils.CreatedResponse(c, "User created successfully", gin.H{
    "user": newUser,
    "token": jwtToken,
})
```

#### `NoContentResponse(c, message)` - 204 No Content
```go
utils.NoContentResponse(c, "User deleted successfully")
```

### Error Responses

#### `ValidationErrorResponse(c, details)` - 400 Bad Request
```go
utils.ValidationErrorResponse(c, "Email is required")
```
**Response:**
```json
{
  "status": "error",
  "message": "Validation failed",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": "Email is required"
  }
}
```

#### `UnauthorizedResponse(c, message)` - 401 Unauthorized
```go
utils.UnauthorizedResponse(c, "Invalid credentials")
```

#### `ForbiddenResponse(c, message)` - 403 Forbidden
```go
utils.ForbiddenResponse(c, "Access denied")
```

#### `NotFoundResponse(c, message)` - 404 Not Found
```go
utils.NotFoundResponse(c, "User not found")
```

#### `ConflictResponse(c, message)` - 409 Conflict
```go
utils.ConflictResponse(c, "User already exists")
```

#### `InternalServerErrorResponse(c, message)` - 500 Internal Server Error
```go
utils.InternalServerErrorResponse(c, "Database connection failed")
```

### Generic Response Functions

#### `SuccessResponse(c, httpStatus, message, data)`
```go
utils.SuccessResponse(c, http.StatusOK, "Custom success", data)
```

#### `ErrorResponse(c, httpStatus, message, errorCode, details)`
```go
utils.ErrorResponse(c, http.StatusBadRequest, "Custom error", "CUSTOM_ERROR", details)
```

#### `FailResponse(c, httpStatus, message, details)`
```go
utils.FailResponse(c, http.StatusBadRequest, "Validation failed", validationErrors)
```

## Usage Examples

### Authentication Controller Examples

#### Registration Success:
```go
func (ac *AuthController) Register(c *gin.Context) {
    // ... validation and business logic ...
    
    utils.CreatedResponse(c, "User registered successfully", gin.H{
        "user":  user,
        "token": token,
    })
}
```

#### Login Failure:
```go
func (ac *AuthController) Login(c *gin.Context) {
    // ... authentication logic ...
    
    if err != nil {
        utils.UnauthorizedResponse(c, "Invalid email or password")
        return
    }
}
```

#### Validation Error:
```go
func (ac *AuthController) Register(c *gin.Context) {
    var req models.UserRegistrationRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ValidationErrorResponse(c, err.Error())
        return
    }
}
```

## Response Examples

### Successful Registration
```json
{
  "status": "success",
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "created_at": "2024-09-26T10:30:00Z",
      "updated_at": "2024-09-26T10:30:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Validation Error
```json
{
  "status": "error",
  "message": "Validation failed",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": "Key: 'UserRegistrationRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"
  }
}
```

### Authentication Error
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

### User Already Exists
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

### Server Error
```json
{
  "status": "error",
  "message": "Failed to register user",
  "error": {
    "code": "INTERNAL_ERROR",
    "details": null
  }
}
```

## HTTP Status Code Mapping

| HTTP Status | Response Function | Status Field | Use Case |
|-------------|------------------|--------------|----------|
| 200 | `OKResponse` | `success` | Successful GET, PUT, PATCH |
| 201 | `CreatedResponse` | `success` | Successful POST (creation) |
| 204 | `NoContentResponse` | `success` | Successful DELETE |
| 400 | `ValidationErrorResponse` | `error` | Bad request, validation errors |
| 401 | `UnauthorizedResponse` | `error` | Authentication required |
| 403 | `ForbiddenResponse` | `error` | Access denied |
| 404 | `NotFoundResponse` | `error` | Resource not found |
| 409 | `ConflictResponse` | `error` | Resource conflict |
| 500 | `InternalServerErrorResponse` | `error` | Server errors |

## Best Practices

### 1. Consistent Messages
Use clear, user-friendly messages:
```go
// Good
utils.UnauthorizedResponse(c, "Invalid email or password")

// Avoid
utils.UnauthorizedResponse(c, "Auth failed")
```

### 2. Meaningful Data Structure
Structure your data objects logically:
```go
utils.OKResponse(c, "Users retrieved successfully", gin.H{
    "users": users,
    "count": len(users),
    "page":  page,
    "total": totalUsers,
})
```

### 3. Error Details
Provide helpful error details for debugging:
```go
if err := c.ShouldBindJSON(&req); err != nil {
    utils.ValidationErrorResponse(c, err.Error())
    return
}
```

### 4. Consistent Error Codes
Use predefined error codes for client-side error handling:
- `VALIDATION_ERROR`: Input validation failed
- `UNAUTHORIZED`: Authentication required
- `FORBIDDEN`: Access denied
- `NOT_FOUND`: Resource not found
- `CONFLICT`: Resource conflict
- `INTERNAL_ERROR`: Server error

## Frontend Integration

### JavaScript/TypeScript Example
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
async function loginUser(email: string, password: string) {
  const response = await fetch('/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  
  const result: APIResponse = await response.json();
  
  if (result.status === 'success') {
    // Handle success
    const { user, token } = result.data;
    localStorage.setItem('token', token);
  } else {
    // Handle error
    console.error(result.message);
    if (result.error?.code === 'VALIDATION_ERROR') {
      // Handle validation errors
    }
  }
}
```

## Migration Guide

If you have existing controllers not using the response wrapper:

### Before:
```go
func GetUser(c *gin.Context) {
    user, err := userService.GetUser(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"user": user})
}
```

### After:
```go
func GetUser(c *gin.Context) {
    user, err := userService.GetUser(id)
    if err != nil {
        utils.NotFoundResponse(c, "User not found")
        return
    }
    utils.OKResponse(c, "User retrieved successfully", gin.H{"user": user})
}
```

This standardized response format makes your API more predictable and easier to consume from frontend applications!

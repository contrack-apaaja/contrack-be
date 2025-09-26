# 📋 Clause Templates API Documentation

## Overview

The Clause Templates API provides CRUD operations and search functionality for managing reusable clause templates in the Contrack system. All endpoints require authentication via JWT token.

## Base URL
```
/api/clauses
```

## Authentication
All endpoints require a valid JWT token in the Authorization header:
```
Authorization: Bearer <jwt-token>
```

## Endpoints

### 1. Create Clause Template
**POST** `/api/clauses`

Creates a new clause template.

#### Request Body
```json
{
  "clause_code": "PAYMENT_001",
  "title": "Standard Payment Terms",
  "type": "Payment",
  "content": "Payment shall be made within 30 days of invoice date...",
  "is_active": true
}
```

#### Field Validations
- `clause_code`: Required, 3-50 characters, must be unique
- `title`: Required, 5-255 characters
- `type`: Required, 3-100 characters
- `content`: Required, minimum 10 characters
- `is_active`: Optional, defaults to `true`

#### Success Response (201)
```json
{
  "status": "success",
  "message": "Clause template created successfully",
  "data": {
    "clause_template": {
      "id": 1,
      "clause_code": "PAYMENT_001",
      "title": "Standard Payment Terms",
      "type": "Payment",
      "content": "Payment shall be made within 30 days of invoice date...",
      "is_active": true,
      "created_at": "2024-09-26T10:30:00Z",
      "updated_at": "2024-09-26T10:30:00Z"
    }
  }
}
```

#### Error Responses
- **400**: Validation error
- **409**: Clause code already exists
- **500**: Internal server error

---

### 2. Get Clause Template by ID
**GET** `/api/clauses/:id`

Retrieves a specific clause template by ID.

#### Path Parameters
- `id`: Integer, clause template ID

#### Success Response (200)
```json
{
  "status": "success",
  "message": "Clause template retrieved successfully",
  "data": {
    "clause_template": {
      "id": 1,
      "clause_code": "PAYMENT_001",
      "title": "Standard Payment Terms",
      "type": "Payment",
      "content": "Payment shall be made within 30 days of invoice date...",
      "is_active": true,
      "created_at": "2024-09-26T10:30:00Z",
      "updated_at": "2024-09-26T10:30:00Z"
    }
  }
}
```

#### Error Responses
- **400**: Invalid ID format
- **404**: Clause template not found
- **500**: Internal server error

---

### 3. Get Clause Template by Code
**GET** `/api/clauses/by-code/:code`

Retrieves a specific clause template by clause code.

#### Path Parameters
- `code`: String, clause code

#### Success Response (200)
Same as Get by ID

#### Error Responses
- **400**: Clause code is required
- **404**: Clause template not found
- **500**: Internal server error

---

### 4. Update Clause Template
**PUT** `/api/clauses/:id`

Updates an existing clause template. All fields are optional.

#### Path Parameters
- `id`: Integer, clause template ID

#### Request Body
```json
{
  "clause_code": "PAYMENT_002",
  "title": "Updated Payment Terms",
  "type": "Payment",
  "content": "Updated payment terms...",
  "is_active": false
}
```

#### Success Response (200)
```json
{
  "status": "success",
  "message": "Clause template updated successfully",
  "data": {
    "clause_template": {
      "id": 1,
      "clause_code": "PAYMENT_002",
      "title": "Updated Payment Terms",
      "type": "Payment",
      "content": "Updated payment terms...",
      "is_active": false,
      "created_at": "2024-09-26T10:30:00Z",
      "updated_at": "2024-09-26T11:00:00Z"
    }
  }
}
```

#### Error Responses
- **400**: Validation error
- **404**: Clause template not found
- **409**: Clause code already exists (if updating code)
- **500**: Internal server error

---

### 5. Delete Clause Template
**DELETE** `/api/clauses/:id`

Deletes a clause template by ID.

#### Path Parameters
- `id`: Integer, clause template ID

#### Success Response (200)
```json
{
  "status": "success",
  "message": "Clause template deleted successfully",
  "data": null
}
```

#### Error Responses
- **400**: Invalid ID format
- **404**: Clause template not found
- **500**: Internal server error

---

### 6. List Clause Templates
**GET** `/api/clauses`

Retrieves a paginated list of clause templates with optional filtering and search.

#### Query Parameters
- `q`: String, search query (searches title and content)
- `type`: String, filter by clause type
- `is_active`: Boolean, filter by active status
- `page`: Integer, page number (default: 1)
- `limit`: Integer, items per page (default: 10, max: 100)
- `sort_by`: String, sort field (`id`, `title`, `type`, `created_at`, `updated_at`)
- `sort_dir`: String, sort direction (`asc`, `desc`)

#### Example Request
```
GET /api/clauses?q=payment&type=Payment&is_active=true&page=1&limit=5&sort_by=title&sort_dir=asc
```

#### Success Response (200)
```json
{
  "status": "success",
  "message": "Clause templates retrieved successfully",
  "data": {
    "clause_templates": [
      {
        "id": 1,
        "clause_code": "PAYMENT_001",
        "title": "Standard Payment Terms",
        "type": "Payment",
        "content": "Payment shall be made within 30 days...",
        "is_active": true,
        "created_at": "2024-09-26T10:30:00Z",
        "updated_at": "2024-09-26T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 5,
      "total": 1,
      "total_pages": 1,
      "has_next": false,
      "has_prev": false
    }
  }
}
```

#### Error Responses
- **400**: Invalid query parameters
- **500**: Internal server error

---

### 7. Search Clause Templates
**GET** `/api/clauses/search`

Performs a focused search on clause templates.

#### Query Parameters
- `q`: String, required, search query
- `limit`: Integer, optional, max results (default: 10)

#### Example Request
```
GET /api/clauses/search?q=payment terms&limit=5
```

#### Success Response (200)
```json
{
  "status": "success",
  "message": "Search completed successfully",
  "data": {
    "query": "payment terms",
    "clause_templates": [...],
    "pagination": {
      "page": 1,
      "limit": 5,
      "total": 3,
      "total_pages": 1,
      "has_next": false,
      "has_prev": false
    }
  }
}
```

#### Error Responses
- **400**: Search query is required
- **500**: Internal server error

---

### 8. Get Clause Types
**GET** `/api/clauses/types`

Retrieves all unique clause types from active clause templates.

#### Success Response (200)
```json
{
  "status": "success",
  "message": "Clause types retrieved successfully",
  "data": {
    "types": [
      "Payment",
      "Delivery",
      "Warranty",
      "Termination",
      "Liability"
    ],
    "count": 5
  }
}
```

#### Error Responses
- **500**: Internal server error

---

### 9. Toggle Clause Template Status
**PATCH** `/api/clauses/:id/toggle-status`

Toggles the active/inactive status of a clause template.

#### Path Parameters
- `id`: Integer, clause template ID

#### Success Response (200)
```json
{
  "status": "success",
  "message": "Clause template activated successfully",
  "data": {
    "clause_template": {
      "id": 1,
      "clause_code": "PAYMENT_001",
      "title": "Standard Payment Terms",
      "type": "Payment",
      "content": "Payment shall be made within 30 days...",
      "is_active": true,
      "created_at": "2024-09-26T10:30:00Z",
      "updated_at": "2024-09-26T11:15:00Z"
    }
  }
}
```

#### Error Responses
- **400**: Invalid ID format
- **404**: Clause template not found
- **500**: Internal server error

---

## Testing Examples

### Using cURL

#### 1. Create a clause template
```bash
curl -X POST http://localhost:8080/api/clauses \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "clause_code": "PAYMENT_001",
    "title": "Standard Payment Terms",
    "type": "Payment",
    "content": "Payment shall be made within 30 days of invoice date. Late payments may incur interest charges.",
    "is_active": true
  }'
```

#### 2. List clause templates with search
```bash
curl "http://localhost:8080/api/clauses?q=payment&page=1&limit=5" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 3. Update a clause template
```bash
curl -X PUT http://localhost:8080/api/clauses/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Payment Terms",
    "content": "Payment shall be made within 45 days of invoice date."
  }'
```

#### 4. Search clause templates
```bash
curl "http://localhost:8080/api/clauses/search?q=payment%20terms&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 5. Get clause types
```bash
curl http://localhost:8080/api/clauses/types \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Using PowerShell

#### Create clause template
```powershell
$headers = @{"Authorization"="Bearer YOUR_JWT_TOKEN"; "Content-Type"="application/json"}
$body = @{
    clause_code = "DELIVERY_001"
    title = "Standard Delivery Terms"
    type = "Delivery"
    content = "Delivery shall be completed within 14 business days."
    is_active = $true
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/clauses" -Method POST -Headers $headers -Body $body
```

#### Search clause templates
```powershell
$headers = @{"Authorization"="Bearer YOUR_JWT_TOKEN"}
Invoke-RestMethod -Uri "http://localhost:8080/api/clauses/search?q=delivery&limit=5" -Headers $headers
```

---

## Database Schema

The clause templates are stored in the `clause_templates` table:

```sql
CREATE TABLE clause_templates (
    id SERIAL PRIMARY KEY,
    clause_code VARCHAR(50) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Indexes
- Primary key on `id`
- Unique index on `clause_code`
- Indexes on `title`, `type`, `is_active`, `created_at`
- Full-text search index on `title` and `content`

---

## Error Handling

All endpoints return standardized error responses:

### Validation Error (400)
```json
{
  "status": "error",
  "message": "Validation failed",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": "Key: 'ClauseTemplateCreateRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag"
  }
}
```

### Not Found Error (404)
```json
{
  "status": "error",
  "message": "Clause template not found",
  "error": {
    "code": "NOT_FOUND",
    "details": null
  }
}
```

### Conflict Error (409)
```json
{
  "status": "error",
  "message": "Clause template with this code already exists",
  "error": {
    "code": "CONFLICT",
    "details": null
  }
}
```

---

## Best Practices

1. **Clause Codes**: Use descriptive, unique codes (e.g., `PAYMENT_001`, `DELIVERY_STD`)
2. **Content**: Write clear, legal language that can be reused across contracts
3. **Types**: Use consistent type names for better organization
4. **Search**: Use the search functionality to avoid creating duplicate clauses
5. **Status Management**: Use the toggle status endpoint to deactivate outdated clauses instead of deleting them

Your Clause Templates API is ready for use! 🚀

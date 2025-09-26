# Contract Management API Testing Guide

## Overview
This guide provides comprehensive testing steps for the Contract Management API implementation.

## Prerequisites
1. Server should be running on port 8080
2. Database should be properly initialized with all tables
3. JWT authentication token required for protected endpoints

## API Endpoints Overview

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user (get JWT token)
- `POST /api/auth/refresh` - Refresh JWT token
- `GET /api/profile` - Get user profile (protected)

### Contract Management
- `POST /api/contracts/` - Create new contract
- `GET /api/contracts/` - List contracts with pagination/filtering
- `GET /api/contracts/:id` - Get specific contract
- `PUT /api/contracts/:id` - Update contract
- `DELETE /api/contracts/:id` - Soft delete contract
- `GET /api/contracts/stats` - Get contract statistics

### Contract Status Management
- `POST /api/contracts/:id/status` - Change contract status
- `GET /api/contracts/:id/status-history` - Get status change history

### Contract Versioning (Separate Group)
- `POST /api/contract-versions/:baseId` - Create new version of contract
- `GET /api/contract-versions/:baseId` - Get all versions of a contract

### Stakeholder Management
- `POST /api/stakeholders/` - Create stakeholder
- `GET /api/stakeholders/` - List stakeholders
- `GET /api/stakeholders/:id` - Get stakeholder details
- `PUT /api/stakeholders/:id` - Update stakeholder
- `DELETE /api/stakeholders/:id` - Delete stakeholder
- `GET /api/stakeholders/types` - Get stakeholder types

## Sample Test Requests

### 1. User Registration & Authentication
```bash
# Register user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'

# Login to get JWT token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 2. Contract Creation
```bash
# Create contract (replace TOKEN with actual JWT)
curl -X POST http://localhost:8080/api/contracts/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "title": "Software Development Agreement",
    "description": "Contract for web application development",
    "contract_type": "SERVICE_AGREEMENT",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-12-31T23:59:59Z",
    "value": 50000.00,
    "currency": "USD",
    "client_info": {
      "company_name": "Tech Corp",
      "contact_person": "John Doe",
      "email": "john@techcorp.com"
    },
    "stakeholder_ids": [],
    "clause_ids": []
  }'
```

### 3. Contract Status Workflow Testing
```bash
# Change status to PENDING_LEGAL_REVIEW
curl -X POST http://localhost:8080/api/contracts/1/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "new_status": "PENDING_LEGAL_REVIEW",
    "reason": "Ready for legal review"
  }'

# Get status history
curl -X GET http://localhost:8080/api/contracts/1/status-history \
  -H "Authorization: Bearer TOKEN"
```

### 4. Contract Search & Filtering
```bash
# List contracts with filters
curl -X GET "http://localhost:8080/api/contracts/?status=DRAFT&page=1&limit=10" \
  -H "Authorization: Bearer TOKEN"

# Get contract statistics
curl -X GET http://localhost:8080/api/contracts/stats \
  -H "Authorization: Bearer TOKEN"
```

### 5. Contract Versioning
```bash
# Create new version of contract
curl -X POST http://localhost:8080/api/contract-versions/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "title": "Software Development Agreement v2",
    "description": "Updated contract with new requirements",
    "changes_summary": "Added new features and updated pricing"
  }'

# Get all versions
curl -X GET http://localhost:8080/api/contract-versions/1 \
  -H "Authorization: Bearer TOKEN"
```

## Expected Status Transitions

### Phase 1 Workflow
1. **DRAFT** → **PENDING_LEGAL_REVIEW** (when ready for review)
2. **PENDING_LEGAL_REVIEW** → **PENDING_SIGNATURE** (after legal approval)
3. **PENDING_SIGNATURE** → **ACTIVE** (when all parties sign)
4. **ACTIVE** → **EXPIRED** (natural expiration)
5. **ACTIVE** → **TERMINATED** (early termination)

### Business Rules Validation
- Only contract owner/creator can modify contracts in DRAFT status
- Status transitions must follow the defined workflow
- Cannot delete contracts in ACTIVE status
- Version numbers increment automatically (CTR-YYYY-MM-XXXXX-VV format)

## Error Handling
- 400: Bad Request (validation errors)
- 401: Unauthorized (missing/invalid JWT)
- 403: Forbidden (access denied)
- 404: Not Found (resource doesn't exist)
- 409: Conflict (invalid status transition)
- 500: Internal Server Error

## Database Tables Created
1. `contracts` - Main contract data with versioning
2. `stakeholders` - Contract stakeholders/parties
3. `contract_stakeholders` - Many-to-many relationship
4. `contract_clauses` - Contract-specific clauses
5. `contract_status_history` - Audit trail for status changes
6. `contract_sequences` - Auto-incrementing contract numbers

## Implementation Notes
- Repository pattern for data access
- Service layer for business logic
- JWT authentication middleware
- Soft deletes with deleted_at timestamps
- JSONB fields for flexible client_info storage
- Full-text search capabilities
- Pagination support
- Comprehensive error handling
- Status transition validation
- Access control based on user roles/ownership

## Next Steps
1. Test all endpoints with actual HTTP requests
2. Verify status transition workflow
3. Test versioning functionality
4. Implement PDF generation for signed contracts
5. Add advanced search capabilities
6. Implement contract templates
7. Add email notifications for status changes
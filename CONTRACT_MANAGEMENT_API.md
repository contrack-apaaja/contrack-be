# Contract Management API Documentation

## Overview

This document provides comprehensive API documentation for the Contract Management System. The system supports CRUD operations for contracts, stakeholders, and contract clauses with proper status workflow management.

## Base URL

```
http://localhost:8080/api
```

## Authentication

All protected endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer <token>
```

## Contract Status Workflow

The contract status follows a strict workflow:

1. **DRAFT** → **PENDING_LEGAL_REVIEW**
2. **PENDING_LEGAL_REVIEW** → **DRAFT** or **PENDING_SIGNATURE**
3. **PENDING_SIGNATURE** → **PENDING_LEGAL_REVIEW** or **ACTIVE**
4. **ACTIVE** → **EXPIRED** or **TERMINATED**
5. **EXPIRED** → **TERMINATED**
6. **TERMINATED** (final status)

## Contract Endpoints

### 1. Create Contract

**POST** `/contracts`

Creates a new contract with initial status `DRAFT`.

**Request Body:**
```json
{
  "project_name": "Highway Construction Project",
  "package_name": "Phase 1 - Road Base",
  "external_reference": "EXT-2024-001",
  "contract_type": "Construction",
  "signing_place": "Jakarta, Indonesia",
  "signing_date": "2024-12-01",
  "total_value": 1500000.00,
  "funding_source": "Government Budget",
  "stakeholders": [
    {
      "stakeholder_id": 1,
      "role_in_contract": "Contractor",
      "representative_name": "John Doe",
      "representative_title": "Project Manager",
      "other_details": {
        "license_number": "CONST-2024-001"
      }
    }
  ],
  "clause_template_ids": [1, 2, 3]
}
```

**Response (201 Created):**
```json
{
  "status": "success",
  "message": "Contract created successfully",
  "data": {
    "id": 1,
    "base_id": "123e4567-e89b-12d3-a456-426614174000",
    "version_number": 1,
    "project_name": "Highway Construction Project",
    "package_name": "Phase 1 - Road Base",
    "contract_number": "CTR-2024-09-00001-V1",
    "external_reference": "EXT-2024-001",
    "contract_type": "Construction",
    "signing_place": "Jakarta, Indonesia",
    "signing_date": "2024-12-01T00:00:00Z",
    "total_value": 1500000.00,
    "funding_source": "Government Budget",
    "status": "DRAFT",
    "created_by": "user-uuid",
    "created_at": "2024-09-26T10:00:00Z",
    "updated_at": "2024-09-26T10:00:00Z",
    "is_deleted": false,
    "stakeholders": [...],
    "clauses": [...]
  }
}
```

### 2. Get Contract

**GET** `/contracts/{id}`

Retrieves a contract by ID with all related data (stakeholders and clauses).

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Contract retrieved successfully",
  "data": {
    "id": 1,
    "base_id": "123e4567-e89b-12d3-a456-426614174000",
    "version_number": 1,
    "project_name": "Highway Construction Project",
    "contract_number": "CTR-2024-09-00001-V1",
    "status": "DRAFT",
    "stakeholders": [
      {
        "id": 1,
        "contract_id": 1,
        "stakeholder_id": 1,
        "role_in_contract": "Contractor",
        "representative_name": "John Doe",
        "representative_title": "Project Manager",
        "stakeholder": {
          "id": 1,
          "legal_name": "ABC Construction Ltd",
          "address": "123 Main St, Jakarta",
          "type": "COMPANY"
        }
      }
    ],
    "clauses": [
      {
        "id": 1,
        "contract_id": 1,
        "clause_template_id": 1,
        "display_order": 1,
        "custom_content": null,
        "clause_template": {
          "id": 1,
          "clause_code": "GENERAL-001",
          "title": "General Terms and Conditions",
          "type": "General",
          "content": "This contract is governed by...",
          "is_active": true
        }
      }
    ]
  }
}
```

### 3. Update Contract

**PUT** `/contracts/{id}`

Updates a contract (only allowed for DRAFT status).

**Request Body:**
```json
{
  "project_name": "Updated Highway Construction Project",
  "total_value": 1600000.00,
  "signing_date": "2024-12-15"
}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Contract updated successfully",
  "data": null
}
```

### 4. Delete Contract

**DELETE** `/contracts/{id}`

Soft deletes a contract (only allowed for DRAFT status).

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Contract deleted successfully",
  "data": null
}
```

### 5. List Contracts

**GET** `/contracts`

Retrieves contracts with search and pagination support.

**Query Parameters:**
- `q` (string): Search query for full-text search
- `status` (string): Filter by contract status
- `contract_type` (string): Filter by contract type
- `funding_source` (string): Filter by funding source
- `signing_date_from` (string): Filter by signing date range (ISO date)
- `signing_date_to` (string): Filter by signing date range (ISO date)
- `value_from` (number): Filter by value range
- `value_to` (number): Filter by value range
- `page` (int): Page number (default: 1)
- `limit` (int): Items per page (default: 10, max: 100)
- `sort_by` (string): Sort field (id, project_name, contract_number, etc.)
- `sort_dir` (string): Sort direction (asc, desc)

**Example Request:**
```
GET /contracts?q=highway&status=DRAFT&page=1&limit=10&sort_by=created_at&sort_dir=desc
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Contracts retrieved successfully",
  "data": {
    "contracts": [
      {
        "id": 1,
        "project_name": "Highway Construction Project",
        "contract_number": "CTR-2024-09-00001-V1",
        "status": "DRAFT",
        "total_value": 1500000.00,
        "created_at": "2024-09-26T10:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 10,
    "pages": 1
  }
}
```

### 6. Change Contract Status

**POST** `/contracts/{id}/status`

Changes the status of a contract with validation of allowed transitions.

**Request Body:**
```json
{
  "status": "PENDING_LEGAL_REVIEW",
  "change_reason": "Ready for legal review",
  "comments": "All documents are complete and ready for review"
}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Contract status changed successfully",
  "data": null
}
```

### 7. Get Status History

**GET** `/contracts/{id}/status-history`

Retrieves the status change history for a contract.

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Status history retrieved successfully",
  "data": [
    {
      "id": 1,
      "contract_id": 1,
      "from_status": "DRAFT",
      "to_status": "PENDING_LEGAL_REVIEW",
      "changed_by": "user-uuid",
      "change_reason": "Ready for legal review",
      "comments": "All documents are complete",
      "changed_at": "2024-09-26T11:00:00Z"
    }
  ]
}
```

### 8. Get Contract Statistics

**GET** `/contracts/stats`

Returns contract statistics for the current user.

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Contract statistics retrieved successfully",
  "data": {
    "DRAFT": 5,
    "PENDING_LEGAL_REVIEW": 3,
    "PENDING_SIGNATURE": 2,
    "ACTIVE": 10,
    "EXPIRED": 1,
    "TERMINATED": 0
  }
}
```

### 9. Create Contract Version

**POST** `/contracts/{baseId}/versions`

Creates a new version of an existing contract.

**Request Body:** Same as Create Contract

**Response (201 Created):**
```json
{
  "status": "success",
  "message": "Contract version created successfully",
  "data": {
    "id": 2,
    "base_id": "123e4567-e89b-12d3-a456-426614174000",
    "version_number": 2,
    "contract_number": "CTR-2024-09-00001-V2",
    "status": "DRAFT"
  }
}
```

### 10. Get Contract Versions

**GET** `/contracts/{baseId}/versions`

Retrieves all versions of a contract.

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Contract versions retrieved successfully",
  "data": [
    {
      "id": 2,
      "version_number": 2,
      "contract_number": "CTR-2024-09-00001-V2",
      "status": "DRAFT",
      "created_at": "2024-09-26T12:00:00Z"
    },
    {
      "id": 1,
      "version_number": 1,
      "contract_number": "CTR-2024-09-00001-V1",
      "status": "ACTIVE",
      "created_at": "2024-09-26T10:00:00Z"
    }
  ]
}
```

## Stakeholder Endpoints

### 1. Create Stakeholder

**POST** `/stakeholders`

**Request Body:**
```json
{
  "legal_name": "ABC Construction Ltd",
  "address": "123 Main Street, Jakarta, Indonesia",
  "type": "COMPANY"
}
```

**Response (201 Created):**
```json
{
  "status": "success",
  "message": "Stakeholder created successfully",
  "data": {
    "id": 1,
    "legal_name": "ABC Construction Ltd",
    "address": "123 Main Street, Jakarta, Indonesia",
    "type": "COMPANY",
    "created_at": "2024-09-26T10:00:00Z",
    "updated_at": "2024-09-26T10:00:00Z",
    "is_deleted": false
  }
}
```

### 2. Get Stakeholder

**GET** `/stakeholders/{id}`

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Stakeholder retrieved successfully",
  "data": {
    "id": 1,
    "legal_name": "ABC Construction Ltd",
    "address": "123 Main Street, Jakarta, Indonesia",
    "type": "COMPANY",
    "created_at": "2024-09-26T10:00:00Z",
    "updated_at": "2024-09-26T10:00:00Z",
    "is_deleted": false
  }
}
```

### 3. Update Stakeholder

**PUT** `/stakeholders/{id}`

**Request Body:**
```json
{
  "legal_name": "ABC Construction Limited",
  "address": "456 New Address, Jakarta, Indonesia"
}
```

### 4. Delete Stakeholder

**DELETE** `/stakeholders/{id}`

### 5. List Stakeholders

**GET** `/stakeholders`

**Query Parameters:**
- `search` (string): Search by legal name
- `type` (string): Filter by stakeholder type
- `page` (int): Page number
- `limit` (int): Items per page

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Stakeholders retrieved successfully",
  "data": {
    "stakeholders": [
      {
        "id": 1,
        "legal_name": "ABC Construction Ltd",
        "type": "COMPANY"
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 10,
    "pages": 1
  }
}
```

### 6. Get Stakeholder Types

**GET** `/stakeholders/types`

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Stakeholder types retrieved successfully",
  "data": [
    "INDIVIDUAL",
    "COMPANY",
    "GOVERNMENT",
    "NGO",
    "OTHER"
  ]
}
```

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "status": "error",
  "message": "Error description",
  "error": {
    "code": "ERROR_CODE",
    "details": "Detailed error information"
  }
}
```

### Common Error Codes

- `VALIDATION_ERROR`: Request validation failed
- `AUTH_ERROR`: Authentication required
- `ACCESS_ERROR`: Access denied
- `NOT_FOUND`: Resource not found
- `STATUS_ERROR`: Invalid status or status transition
- `CREATE_ERROR`: Failed to create resource
- `UPDATE_ERROR`: Failed to update resource
- `DELETE_ERROR`: Failed to delete resource
- `FETCH_ERROR`: Failed to retrieve resource

## Business Rules

### Contract Editing Rules
- Only contracts with status `DRAFT` can be edited or deleted
- Contract numbers are automatically generated with format: `CTR-YYYY-MM-XXXXX-VV`
- Signing date cannot be in the future
- Total value must be positive
- Status transitions must follow the defined workflow

### Access Control Rules (Phase 1)
- Users can only access contracts they created
- All contract operations are restricted to the contract creator
- Future phases will implement role-based access control

### Versioning Rules
- New contract versions always start with status `DRAFT`
- Version numbers increment automatically
- Contract numbers include version suffix (e.g., `-V1`, `-V2`)
- All versions share the same `base_id`

## Testing Examples

### Create a Complete Contract

```bash
# 1. Create stakeholder
curl -X POST http://localhost:8080/api/stakeholders \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "legal_name": "ABC Construction Ltd",
    "address": "123 Main Street, Jakarta",
    "type": "COMPANY"
  }'

# 2. Create contract with stakeholder
curl -X POST http://localhost:8080/api/contracts \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "project_name": "Highway Construction Project",
    "contract_type": "Construction",
    "total_value": 1500000.00,
    "stakeholders": [
      {
        "stakeholder_id": 1,
        "role_in_contract": "Contractor",
        "representative_name": "John Doe",
        "representative_title": "Project Manager"
      }
    ],
    "clause_template_ids": [1, 2, 3]
  }'

# 3. Submit for legal review
curl -X POST http://localhost:8080/api/contracts/1/status \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "PENDING_LEGAL_REVIEW",
    "change_reason": "Ready for legal review",
    "comments": "All documents completed"
  }'
```

This completes the comprehensive Contract Management API implementation for Phase 1!
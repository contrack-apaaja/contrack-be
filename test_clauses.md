# 🧪 Clause Templates API Testing Guide

## Prerequisites

1. **Start the server**:
```bash
go run ./cmd/server
```

2. **Get JWT token** by logging in:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"your@email.com","password":"yourpassword"}'
```

**Save the token** from the response for the following tests.

## Test Sequence

### 1. Create Sample Clause Templates

#### Payment Clause
```bash
curl -X POST http://localhost:8080/api/clauses \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "clause_code": "PAYMENT_001",
    "title": "Standard Payment Terms",
    "type": "Payment",
    "content": "Payment shall be made within 30 days of invoice date. Late payments may incur interest charges at 1.5% per month.",
    "is_active": true
  }'
```

#### Delivery Clause
```bash
curl -X POST http://localhost:8080/api/clauses \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "clause_code": "DELIVERY_001",
    "title": "Standard Delivery Terms",
    "type": "Delivery",
    "content": "Delivery shall be completed within 14 business days from order confirmation. Risk of loss transfers upon delivery.",
    "is_active": true
  }'
```

#### Warranty Clause
```bash
curl -X POST http://localhost:8080/api/clauses \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "clause_code": "WARRANTY_001",
    "title": "Limited Warranty",
    "type": "Warranty",
    "content": "Seller warrants that goods will be free from defects in materials and workmanship for a period of 12 months from delivery.",
    "is_active": true
  }'
```

### 2. Test Retrieval

#### Get clause by ID
```bash
curl -X GET http://localhost:8080/api/clauses/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get clause by code
```bash
curl -X GET http://localhost:8080/api/clauses/by-code/PAYMENT_001 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3. Test Listing and Search

#### List all clauses
```bash
curl -X GET http://localhost:8080/api/clauses \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### List with pagination
```bash
curl -X GET "http://localhost:8080/api/clauses?page=1&limit=2&sort_by=title&sort_dir=asc" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Search by query
```bash
curl -X GET "http://localhost:8080/api/clauses/search?q=payment&limit=5" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Filter by type
```bash
curl -X GET "http://localhost:8080/api/clauses?type=Payment&is_active=true" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. Test Updates

#### Update clause content
```bash
curl -X PUT http://localhost:8080/api/clauses/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Payment shall be made within 45 days of invoice date. Late payments may incur interest charges at 2% per month."
  }'
```

#### Toggle clause status
```bash
curl -X PATCH http://localhost:8080/api/clauses/1/toggle-status \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 5. Test Utility Endpoints

#### Get all clause types
```bash
curl -X GET http://localhost:8080/api/clauses/types \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 6. Test Error Cases

#### Try to create duplicate clause code
```bash
curl -X POST http://localhost:8080/api/clauses \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "clause_code": "PAYMENT_001",
    "title": "Duplicate Payment Terms",
    "type": "Payment",
    "content": "This should fail due to duplicate code."
  }'
```

#### Try to get non-existent clause
```bash
curl -X GET http://localhost:8080/api/clauses/999 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Try invalid validation
```bash
curl -X POST http://localhost:8080/api/clauses \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "clause_code": "AB",
    "title": "Too",
    "type": "X",
    "content": "Short"
  }'
```

### 7. Clean Up (Optional)

#### Delete test clauses
```bash
curl -X DELETE http://localhost:8080/api/clauses/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

curl -X DELETE http://localhost:8080/api/clauses/2 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

curl -X DELETE http://localhost:8080/api/clauses/3 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Expected Results

✅ **All operations should work correctly and return proper response formats**
✅ **Search should find clauses by title and content**
✅ **Pagination should work with proper metadata**
✅ **Validation should prevent invalid data**
✅ **Unique constraints should be enforced**
✅ **Status toggling should work**

## PowerShell Alternative

If you're using PowerShell, here's an example for creating a clause:

```powershell
$headers = @{
    "Authorization" = "Bearer YOUR_JWT_TOKEN"
    "Content-Type" = "application/json"
}

$body = @{
    clause_code = "TERMINATION_001"
    title = "Standard Termination Clause"
    type = "Termination"
    content = "Either party may terminate this agreement with 30 days written notice."
    is_active = $true
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/clauses" -Method POST -Headers $headers -Body $body
```

Your Clause Templates API is ready for testing! 🎯

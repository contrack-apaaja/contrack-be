#!/bin/bash

# Test script untuk AI contract analysis
BASE_URL="http://localhost:8080/api"

echo "=== Testing AI Contract Analysis ==="

# 1. Login user
echo "1. Login user..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')

echo "Login response: $LOGIN_RESPONSE"

# Extract token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Token: $TOKEN"

if [ -z "$TOKEN" ]; then
  echo "Error: Could not extract token from login response"
  exit 1
fi

# 2. Test contract analysis endpoint
echo "2. Testing contract analysis endpoint..."
CONTRACT_ANALYSIS_RESPONSE=$(curl -s -X POST "$BASE_URL/ai/analyze-contract" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "contract_id": 1,
    "clause_template_ids": [1, 2, 3]
  }')

echo "Contract analysis response: $CONTRACT_ANALYSIS_RESPONSE"

# 3. Test old clause analysis endpoint (for comparison)
echo "3. Testing old clause analysis endpoint..."
CLAUSE_ANALYSIS_RESPONSE=$(curl -s -X POST "$BASE_URL/ai/analyze" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "clause_id": 1
  }')

echo "Clause analysis response: $CLAUSE_ANALYSIS_RESPONSE"

echo "=== AI Contract Analysis Test Complete ==="

#!/bin/bash

# Full workflow test: Create analysis then get recommendations
BASE_URL="http://localhost:8080/api"

echo "=== Full Workflow Test: Analysis + Recommendations ==="

# 1. Login to get token
echo "1. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}')

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ Error: Could not get token"
  exit 1
fi

echo "✅ Token obtained: ${TOKEN:0:20}..."

# 2. Create contract analysis first
echo "2. Creating contract analysis..."
ANALYSIS_RESPONSE=$(curl -s -X POST "${BASE_URL}/ai/analyze-contract" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"contract_id": 1, "clause_template_ids": [1, 2, 3]}')

echo "Analysis response: $ANALYSIS_RESPONSE"

# Check if analysis was successful
if echo "$ANALYSIS_RESPONSE" | grep -q "success"; then
  echo "✅ Contract analysis created successfully"
else
  echo "❌ Error creating analysis"
  echo "Response: $ANALYSIS_RESPONSE"
  exit 1
fi

# 3. Now test getting recommendations
echo "3. Getting contract recommendations..."
RECOMMENDATIONS_RESPONSE=$(curl -s -X GET "${BASE_URL}/ai/contract/1/recommendations" \
  -H "Authorization: Bearer $TOKEN")

echo "Recommendations response: $RECOMMENDATIONS_RESPONSE"

# Check if recommendations were retrieved successfully
if echo "$RECOMMENDATIONS_RESPONSE" | grep -q "success"; then
  echo "✅ Contract recommendations retrieved successfully"
else
  echo "❌ Error getting recommendations"
  echo "Response: $RECOMMENDATIONS_RESPONSE"
fi

echo "=== Full Workflow Test Complete ==="

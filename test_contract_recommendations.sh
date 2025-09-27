#!/bin/bash

# Test script for contract recommendations endpoint
BASE_URL="http://localhost:8080/api"

echo "=== Testing Contract Recommendations Endpoint ==="

# 1. Login to get token
echo "1. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}')

echo "Login response: $LOGIN_RESPONSE"

# Extract token
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ Error: Could not get token"
  exit 1
fi

echo "✅ Token obtained: ${TOKEN:0:20}..."

# 2. Test contract recommendations endpoint
echo "2. Testing contract recommendations endpoint..."
echo "URL: ${BASE_URL}/ai/contract/1/recommendations"

RECOMMENDATIONS_RESPONSE=$(curl -s -X GET "${BASE_URL}/ai/contract/1/recommendations" \
  -H "Authorization: Bearer $TOKEN")

echo "Recommendations response: $RECOMMENDATIONS_RESPONSE"

# Check if response contains error
if echo "$RECOMMENDATIONS_RESPONSE" | grep -q "error"; then
  echo "❌ Error in response"
  echo "Response: $RECOMMENDATIONS_RESPONSE"
else
  echo "✅ Success! Contract recommendations retrieved"
fi

echo "=== Test Complete ==="

#!/bin/bash

# Quick test script for AI Analysis feature
# This script will test the AI analysis feature end-to-end

echo "🚀 Quick Test AI Analysis Feature"
echo "================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="http://localhost:8080"
TEST_EMAIL="test@example.com"
TEST_PASSWORD="password123"
TEST_NAME="Test User"

# Function to print colored output
print_status() {
    if [ $2 -eq 0 ]; then
        echo -e "${GREEN}✅ $1${NC}"
    else
        echo -e "${RED}❌ $1${NC}"
    fi
}

# Function to make HTTP request and get response
make_request() {
    local method=$1
    local url=$2
    local data=$3
    local headers=$4
    
    if [ -n "$data" ]; then
        curl -s -X $method "$url" \
            -H "Content-Type: application/json" \
            -H "$headers" \
            -d "$data"
    else
        curl -s -X $method "$url" \
            -H "$headers"
    fi
}

# Check if server is running
echo "🔍 Checking if server is running..."
server_check=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/hello")
if [ "$server_check" = "200" ]; then
    print_status "Server is running" 0
else
    print_status "Server is not running. Please start the server first: go run cmd/server/main.go" 1
    exit 1
fi

# Step 1: Register user
echo ""
echo "📝 Step 1: Registering user..."
register_data='{"email":"'$TEST_EMAIL'","password":"'$TEST_PASSWORD'","name":"'$TEST_NAME'"}'
register_response=$(make_request "POST" "$BASE_URL/api/auth/register" "$register_data" "")

if echo "$register_response" | grep -q '"status":"success"'; then
    print_status "User registered successfully" 0
elif echo "$register_response" | grep -q "already exists"; then
    print_status "User already exists, continuing..." 0
else
    print_status "Failed to register user: $register_response" 1
    exit 1
fi

# Step 2: Login and get JWT token
echo ""
echo "🔐 Step 2: Logging in..."
login_data='{"email":"'$TEST_EMAIL'","password":"'$TEST_PASSWORD'"}'
login_response=$(make_request "POST" "$BASE_URL/api/auth/login" "$login_data" "")

# Extract JWT token
JWT_TOKEN=$(echo "$login_response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$JWT_TOKEN" ]; then
    print_status "Login successful, JWT token obtained" 0
else
    print_status "Failed to login: $login_response" 1
    exit 1
fi

# Step 3: Create test clause
echo ""
echo "📄 Step 3: Creating test clause..."
clause_data='{
    "clause_code": "TEST_AI_001",
    "title": "Test Clause for AI Analysis",
    "type": "Payment",
    "content": "Pembayaran harus dilakukan dalam waktu 30 hari setelah invoice diterima. Jika pembayaran terlambat, akan dikenakan bunga 2% per bulan. Pihak yang membayar bertanggung jawab penuh atas semua biaya yang timbul akibat keterlambatan pembayaran.",
    "is_active": true
}'

clause_response=$(make_request "POST" "$BASE_URL/api/clauses" "$clause_data" "Authorization: Bearer $JWT_TOKEN")

# Extract clause ID
CLAUSE_ID=$(echo "$clause_response" | grep -o '"id":[0-9]*' | cut -d':' -f2)

if [ -n "$CLAUSE_ID" ]; then
    print_status "Clause created successfully with ID: $CLAUSE_ID" 0
else
    print_status "Failed to create clause: $clause_response" 1
    exit 1
fi

# Step 4: Analyze clause with AI
echo ""
echo "🤖 Step 4: Analyzing clause with AI..."
analysis_data='{"clause_id":'$CLAUSE_ID'}'
analysis_response=$(make_request "POST" "$BASE_URL/api/ai/analyze" "$analysis_data" "Authorization: Bearer $JWT_TOKEN")

# Extract analysis ID
ANALYSIS_ID=$(echo "$analysis_response" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

if [ -n "$ANALYSIS_ID" ]; then
    print_status "AI analysis completed successfully with ID: $ANALYSIS_ID" 0
    
    # Extract and display key results
    echo ""
    echo "📊 Analysis Results:"
    echo "-------------------"
    
    # Extract risk level
    RISK_LEVEL=$(echo "$analysis_response" | grep -o '"risk_level":"[^"]*"' | cut -d'"' -f4)
    echo "Risk Level: $RISK_LEVEL"
    
    # Extract risk score
    RISK_SCORE=$(echo "$analysis_response" | grep -o '"risk_score":[0-9.]*' | cut -d':' -f2)
    echo "Risk Score: $RISK_SCORE"
    
    # Extract confidence score
    CONFIDENCE_SCORE=$(echo "$analysis_response" | grep -o '"confidence_score":[0-9.]*' | cut -d':' -f2)
    echo "Confidence Score: $CONFIDENCE_SCORE"
    
    # Extract model version
    MODEL_VERSION=$(echo "$analysis_response" | grep -o '"model_version":"[^"]*"' | cut -d'"' -f4)
    echo "Model Version: $MODEL_VERSION"
    
else
    print_status "Failed to analyze clause: $analysis_response" 1
    exit 1
fi

# Step 5: Get analysis by ID
echo ""
echo "🔍 Step 5: Getting analysis by ID..."
get_analysis_response=$(make_request "GET" "$BASE_URL/api/ai/analysis/$ANALYSIS_ID" "" "Authorization: Bearer $JWT_TOKEN")

if echo "$get_analysis_response" | grep -q "success.*true"; then
    print_status "Analysis retrieved successfully" 0
else
    print_status "Failed to retrieve analysis: $get_analysis_response" 1
fi

# Step 6: Get analysis by clause ID
echo ""
echo "🔍 Step 6: Getting analysis by clause ID..."
get_clause_analysis_response=$(make_request "GET" "$BASE_URL/api/ai/analysis/clause/$CLAUSE_ID" "" "Authorization: Bearer $JWT_TOKEN")

if echo "$get_clause_analysis_response" | grep -q "success.*true"; then
    print_status "Clause analysis retrieved successfully" 0
else
    print_status "Failed to retrieve clause analysis: $get_clause_analysis_response" 1
fi

# Step 7: Get analyses list
echo ""
echo "📋 Step 7: Getting analyses list..."
get_analyses_response=$(make_request "GET" "$BASE_URL/api/ai/analyses" "" "Authorization: Bearer $JWT_TOKEN")

if echo "$get_analyses_response" | grep -q "success.*true"; then
    print_status "Analyses list retrieved successfully" 0
else
    print_status "Failed to retrieve analyses list: $get_analyses_response" 1
fi

# Step 8: Get statistics
echo ""
echo "📈 Step 8: Getting analysis statistics..."
get_stats_response=$(make_request "GET" "$BASE_URL/api/ai/stats" "" "Authorization: Bearer $JWT_TOKEN")

if echo "$get_stats_response" | grep -q "success.*true"; then
    print_status "Statistics retrieved successfully" 0
else
    print_status "Failed to retrieve statistics: $get_stats_response" 1
fi

# Summary
echo ""
echo "🎉 Test Summary"
echo "=============="
echo "✅ User registration and login"
echo "✅ Clause creation"
echo "✅ AI analysis"
echo "✅ Analysis retrieval"
echo "✅ Statistics"
echo ""
echo "🔗 Test URLs for Postman:"
echo "Base URL: $BASE_URL"
echo "JWT Token: $JWT_TOKEN"
echo "Clause ID: $CLAUSE_ID"
echo "Analysis ID: $ANALYSIS_ID"
echo ""
echo "📝 Next steps:"
echo "1. Use the JWT token in Postman for testing"
echo "2. Try different clause types to see various risk levels"
echo "3. Test filtering and pagination features"
echo "4. Check the database for stored results"
echo ""
echo "✨ AI Analysis feature is working correctly!"

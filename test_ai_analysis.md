# Test AI Analysis Feature

## Setup Environment

1. Set environment variable:
```bash
export OPENAI_API_KEY=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
```

2. Run database migration:
```bash
# Pastikan migration 003_create_clause_risk_analyses.sql sudah dijalankan
```

3. Start server:
```bash
go run cmd/server/main.go
```

## Test Cases

### 1. Create Test Clause

```bash
curl -X POST http://localhost:8080/api/clauses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "clause_code": "TEST_AI_001",
    "title": "Klausul Test AI",
    "type": "Payment",
    "content": "Pembayaran harus dilakukan dalam waktu 30 hari setelah invoice diterima. Jika pembayaran terlambat, akan dikenakan bunga 2% per bulan. Pihak yang membayar bertanggung jawab penuh atas semua biaya yang timbul akibat keterlambatan pembayaran.",
    "is_active": true
  }'
```

### 2. Analyze Clause Risk

```bash
curl -X POST http://localhost:8080/api/ai/analyze \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "clause_id": 1
  }'
```

Expected response:
```json
{
  "success": true,
  "message": "Analysis completed successfully",
  "data": {
    "analysis": {
      "id": 1,
      "clause_id": 1,
      "risk_level": "medium",
      "risk_score": 65.5,
      "analysis_summary": "Klausul pembayaran ini memiliki risiko sedang...",
      "identified_risks": [
        "Bunga 2% per bulan mungkin terlalu tinggi",
        "Tanggung jawab penuh dapat menimbulkan ketidakadilan",
        "Tidak ada klausul force majeure"
      ],
      "recommendations": [
        "Pertimbangkan untuk mengurangi tingkat bunga",
        "Tambahkan klausul force majeure",
        "Sertakan mekanisme penyelesaian sengketa"
      ],
      "legal_implications": "Klausul ini dapat melanggar prinsip keseimbangan kontrak...",
      "compliance_notes": "Perlu memastikan kepatuhan terhadap UU Perlindungan Konsumen...",
      "confidence_score": 85.0,
      "model_version": "gpt-3.5-turbo-v1.0",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    "clause": {
      "id": 1,
      "clause_code": "TEST_AI_001",
      "title": "Klausul Test AI",
      "type": "Payment",
      "content": "Pembayaran harus dilakukan dalam waktu 30 hari...",
      "is_active": true,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  }
}
```

### 3. Get Analysis by ID

```bash
curl -X GET http://localhost:8080/api/ai/analysis/1 \
  -H "Authorization: Bearer <your-jwt-token>"
```

### 4. Get Analysis by Clause ID

```bash
curl -X GET http://localhost:8080/api/ai/analysis/clause/1 \
  -H "Authorization: Bearer <your-jwt-token>"
```

### 5. Get All Analyses

```bash
curl -X GET http://localhost:8080/api/ai/analyses \
  -H "Authorization: Bearer <your-jwt-token>"
```

### 6. Filter High Risk Analyses

```bash
curl -X GET "http://localhost:8080/api/ai/analyses?risk_level=high&min_risk_score=70" \
  -H "Authorization: Bearer <your-jwt-token>"
```

### 7. Get Analysis Statistics

```bash
curl -X GET http://localhost:8080/api/ai/stats \
  -H "Authorization: Bearer <your-jwt-token>"
```

### 8. Delete Analysis

```bash
curl -X DELETE http://localhost:8080/api/ai/analysis/1 \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Test Scenarios

### Scenario 1: Low Risk Clause
```json
{
  "clause_code": "LOW_RISK_001",
  "title": "Klausul Risiko Rendah",
  "type": "General",
  "content": "Para pihak sepakat untuk menyelesaikan sengketa melalui mediasi terlebih dahulu sebelum mengajukan gugatan ke pengadilan."
}
```

Expected: Risk level "low", risk score < 25

### Scenario 2: High Risk Clause
```json
{
  "clause_code": "HIGH_RISK_001",
  "title": "Klausul Risiko Tinggi",
  "type": "Liability",
  "content": "Pihak pertama tidak bertanggung jawab atas kerugian apapun yang timbul dari kontrak ini, termasuk kerugian langsung, tidak langsung, konsekuensial, atau punitif, tanpa batas waktu dan tanpa pengecualian apapun."
}
```

Expected: Risk level "high" or "critical", risk score > 50

### Scenario 3: Medium Risk Clause
```json
{
  "clause_code": "MEDIUM_RISK_001",
  "title": "Klausul Risiko Sedang",
  "type": "Payment",
  "content": "Pembayaran dilakukan dalam 14 hari setelah barang diterima. Keterlambatan pembayaran akan dikenakan denda 1% per hari."
}
```

Expected: Risk level "medium", risk score 25-50

## Error Testing

### 1. Invalid Clause ID
```bash
curl -X POST http://localhost:8080/api/ai/analyze \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "clause_id": 999
  }'
```

Expected: 404 Not Found

### 2. Missing OpenAI API Key
```bash
# Unset API key
unset OPENAI_API_KEY

curl -X POST http://localhost:8080/api/ai/analyze \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "clause_id": 1
  }'
```

Expected: 500 Internal Server Error

### 3. Invalid Request Format
```bash
curl -X POST http://localhost:8080/api/ai/analyze \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "invalid_field": "test"
  }'
```

Expected: 400 Bad Request

## Performance Testing

### 1. Concurrent Analysis Requests
```bash
# Run multiple analyses simultaneously
for i in {1..5}; do
  curl -X POST http://localhost:8080/api/ai/analyze \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <your-jwt-token>" \
    -d "{\"clause_id\": $i}" &
done
wait
```

### 2. Large Clause Content
```bash
# Test with very long clause content
curl -X POST http://localhost:8080/api/ai/analyze \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "clause_id": 1
  }'
```

## Monitoring

### 1. Check Logs
```bash
# Monitor server logs for AI analysis requests
tail -f server.log | grep "AI analysis"
```

### 2. Database Queries
```sql
-- Check analysis count
SELECT COUNT(*) FROM clause_risk_analyses;

-- Check risk distribution
SELECT risk_level, COUNT(*) FROM clause_risk_analyses GROUP BY risk_level;

-- Check average scores
SELECT AVG(risk_score), AVG(confidence_score) FROM clause_risk_analyses;
```

### 3. OpenAI Usage
- Monitor OpenAI API usage in dashboard
- Check token consumption
- Monitor rate limits

## Troubleshooting

### Common Issues

1. **OpenAI API Error 401**
   - Check API key validity
   - Verify API key has sufficient credits

2. **Database Connection Error**
   - Check database connection
   - Verify migration has been run

3. **Analysis Timeout**
   - Check OpenAI API response time
   - Verify network connectivity

4. **Invalid JSON Response**
   - Check OpenAI model response format
   - Verify prompt engineering

### Debug Commands

```bash
# Check environment variables
env | grep OPENAI

# Test database connection
psql $DATABASE_URL -c "SELECT 1;"

# Check server status
curl http://localhost:8080/api/hello

# Test authentication
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password"}'
```

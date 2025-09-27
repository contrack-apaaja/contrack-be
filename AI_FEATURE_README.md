# AI Analysis Feature

## Overview
Fitur AI rekomendasi untuk analisis risiko klausul kontrak menggunakan OpenAI GPT-3.5-turbo. Sistem ini menganalisis metadata klausul dan memberikan penilaian risiko beserta rekomendasi.

## Features
- ✅ Analisis risiko klausul kontrak menggunakan AI
- ✅ Penilaian risiko dengan skor 0-100
- ✅ Kategorisasi risiko (low, medium, high, critical)
- ✅ Identifikasi risiko spesifik
- ✅ Rekomendasi perbaikan
- ✅ Analisis implikasi hukum
- ✅ Catatan kepatuhan
- ✅ Skor kepercayaan analisis
- ✅ Penyimpanan hasil analisis
- ✅ API endpoints lengkap
- ✅ Pagination dan filtering
- ✅ Statistik analisis

## Quick Start

### 1. Setup Environment
```bash
# Set OpenAI API key
export OPENAI_API_KEY=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

# Or add to .env file
echo "OPENAI_API_KEY=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" >> .env
```

### 2. Run Database Migration
```bash
# Run migration script
./run_migration.sh

# Or manually
psql "$DATABASE_URL" -f migrations/003_create_clause_risk_analyses.sql
```

### 3. Start Server
```bash
go run cmd/server/main.go
```

### 4. Test the Feature
```bash
# Create a test clause
curl -X POST http://localhost:8080/api/clauses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "clause_code": "TEST_001",
    "title": "Test Clause",
    "type": "Payment",
    "content": "Pembayaran harus dilakukan dalam waktu 30 hari setelah invoice diterima.",
    "is_active": true
  }'

# Analyze the clause
curl -X POST http://localhost:8080/api/ai/analyze \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{"clause_id": 1}'
```

## API Endpoints

### Core Endpoints
- `POST /api/ai/analyze` - Analyze clause risk
- `GET /api/ai/analysis/:id` - Get analysis by ID
- `GET /api/ai/analysis/clause/:clause_id` - Get analysis by clause ID
- `GET /api/ai/analyses` - Get analyses with pagination
- `DELETE /api/ai/analysis/:id` - Delete analysis
- `GET /api/ai/stats` - Get analysis statistics

### Query Parameters for `/api/ai/analyses`
- `clause_id` - Filter by clause ID
- `risk_level` - Filter by risk level (low, medium, high, critical)
- `min_risk_score` - Minimum risk score (0-100)
- `max_risk_score` - Maximum risk score (0-100)
- `min_confidence` - Minimum confidence score (0-100)
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)
- `sort_by` - Sort field (id, clause_id, risk_level, risk_score, created_at, updated_at)
- `sort_dir` - Sort direction (asc, desc)

## Risk Levels

| Level | Score Range | Description |
|-------|-------------|-------------|
| Low | 0-24 | Risiko minimal, klausul aman |
| Medium | 25-49 | Risiko sedang, perlu perhatian |
| High | 50-74 | Risiko tinggi, perlu revisi |
| Critical | 75-100 | Risiko kritis, tidak disarankan |

## Database Schema

### clause_risk_analyses Table
```sql
CREATE TABLE clause_risk_analyses (
    id SERIAL PRIMARY KEY,
    clause_id INTEGER NOT NULL REFERENCES clause_templates(id) ON DELETE CASCADE,
    risk_level VARCHAR(20) NOT NULL CHECK (risk_level IN ('low', 'medium', 'high', 'critical')),
    risk_score DECIMAL(5,2) NOT NULL CHECK (risk_score >= 0 AND risk_score <= 100),
    analysis_summary TEXT NOT NULL,
    identified_risks JSONB NOT NULL DEFAULT '[]',
    recommendations JSONB NOT NULL DEFAULT '[]',
    legal_implications TEXT NOT NULL,
    compliance_notes TEXT NOT NULL,
    confidence_score DECIMAL(5,2) NOT NULL CHECK (confidence_score >= 0 AND confidence_score <= 100),
    model_version VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## File Structure

```
internal/
├── models/
│   └── ai_analysis.go          # AI analysis models
├── services/
│   └── ai/
│       └── ai.go               # OpenAI service
├── repository/
│   └── ai_repo.go              # AI analysis repository
├── controllers/
│   └── ai.go                   # AI analysis controller
└── router/
    └── router.go               # Updated with AI routes

migrations/
└── 003_create_clause_risk_analyses.sql

docs/
├── AI_ANALYSIS_API_DOCUMENTATION.md
├── test_ai_analysis.md
└── AI_FEATURE_README.md
```

## Configuration

### Environment Variables
```bash
# Required
OPENAI_API_KEY=your-openai-api-key

# Optional (with defaults)
PORT=8080
DATABASE_URL=your-database-url
JWT_SECRET=your-jwt-secret
```

### OpenAI Configuration
- Model: GPT-3.5-turbo
- Max Tokens: 2000
- Temperature: 0.3
- Timeout: 60 seconds

## Usage Examples

### 1. Basic Analysis
```bash
curl -X POST http://localhost:8080/api/ai/analyze \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"clause_id": 1}'
```

### 2. Filter High Risk Analyses
```bash
curl -X GET "http://localhost:8080/api/ai/analyses?risk_level=high&min_risk_score=70" \
  -H "Authorization: Bearer <token>"
```

### 3. Get Analysis Statistics
```bash
curl -X GET http://localhost:8080/api/ai/stats \
  -H "Authorization: Bearer <token>"
```

## Error Handling

### Common Errors
- **401 Unauthorized**: Invalid or missing JWT token
- **404 Not Found**: Clause or analysis not found
- **400 Bad Request**: Invalid request format
- **500 Internal Server Error**: OpenAI API error or database error

### Error Response Format
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

## Performance Considerations

### Rate Limiting
- OpenAI API: 3 requests per minute (free tier)
- Application: 10 requests per minute per user
- Database: Optimized with indexes

### Optimization Tips
1. Cache analysis results for identical clauses
2. Batch multiple analyses when possible
3. Monitor OpenAI API usage and costs
4. Use pagination for large result sets
5. Implement proper error handling and retries

## Security

### Authentication
- All endpoints require JWT authentication
- Token validation on every request
- User-specific access control

### Data Protection
- Sensitive data encrypted in transit
- API keys stored securely
- Input validation and sanitization
- SQL injection prevention

## Monitoring

### Logs
- AI analysis requests and responses
- OpenAI API usage and errors
- Database query performance
- Authentication events

### Metrics
- Analysis success rate
- Average response time
- Risk level distribution
- API usage statistics

## Troubleshooting

### Common Issues

1. **OpenAI API Error 401**
   - Check API key validity
   - Verify API key has sufficient credits

2. **Database Connection Error**
   - Check database connection string
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
psql "$DATABASE_URL" -c "SELECT 1;"

# Check server status
curl http://localhost:8080/api/hello

# Test authentication
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password"}'
```

## Contributing

### Development Setup
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

### Testing
```bash
# Run tests
go test ./...

# Test AI feature specifically
go test ./internal/services/ai/...
go test ./internal/controllers/...
go test ./internal/repository/...
```

### Code Style
- Follow Go conventions
- Use meaningful variable names
- Add comments for complex logic
- Write comprehensive tests

## License
This project is licensed under the MIT License.

## Support
For support and questions:
- Create an issue in the repository
- Check the documentation
- Review the troubleshooting guide

## Changelog

### v1.0.0 (2024-01-15)
- Initial release of AI analysis feature
- OpenAI GPT-3.5-turbo integration
- Complete API endpoints
- Database schema and migrations
- Comprehensive documentation
- Test cases and examples

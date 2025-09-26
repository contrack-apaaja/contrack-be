#!/bin/bash

# Script to run database migration for AI analysis feature

echo "🔄 Running database migration for AI analysis feature..."

# Load environment variables from .env file
if [ -f .env ]; then
    echo "📄 Loading environment variables from .env file..."
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "⚠️  Warning: .env file not found"
fi

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "❌ Error: DATABASE_URL environment variable is not set"
    echo "Please set DATABASE_URL in your .env file or environment"
    exit 1
fi

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo "❌ Error: psql command not found"
    echo "Please install PostgreSQL client tools"
    exit 1
fi

# Run the migration
echo "📊 Creating clause_risk_analyses table..."
psql "$DATABASE_URL" -f migrations/003_create_clause_risk_analyses.sql

if [ $? -eq 0 ]; then
    echo "✅ Migration completed successfully!"
    echo "🎉 AI analysis feature is now ready to use"
else
    echo "❌ Migration failed!"
    exit 1
fi

echo ""
echo "📋 Next steps:"
echo "1. Set OPENAI_API_KEY in your environment"
echo "2. Start the server: go run cmd/server/main.go"
echo "3. Test the AI analysis feature using the test_ai_analysis.md guide"

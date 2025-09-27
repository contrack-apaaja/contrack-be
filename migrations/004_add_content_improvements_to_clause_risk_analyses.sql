-- Add content_improvements field to clause_risk_analyses table
-- This field will store AI-generated content improvement recommendations

ALTER TABLE clause_risk_analyses 
ADD COLUMN content_improvements JSONB NOT NULL DEFAULT '[]';

-- Add comment for documentation
COMMENT ON COLUMN clause_risk_analyses.content_improvements IS 'JSON array of AI-generated content improvement recommendations for the clause';

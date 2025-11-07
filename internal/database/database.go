package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewConnection(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create or alter candidates table to include all required fields
	if err := createOrAlterCandidatesTable(db); err != nil {
		return nil, fmt.Errorf("failed to setup candidates table: %w", err)
	}

	return db, nil
}

func createOrAlterCandidatesTable(db *sql.DB) error {
	// Drop unique constraint on email if it exists (to allow NULL values)
	_, _ = db.Exec(`ALTER TABLE candidates DROP CONSTRAINT IF EXISTS candidates_email_key`)

	// Add missing columns if they don't exist
	alterTableQueries := []string{
		`ALTER TABLE candidates ADD COLUMN IF NOT EXISTS position VARCHAR(255)`,
		`ALTER TABLE candidates ADD COLUMN IF NOT EXISTS experience VARCHAR(50)`,
		`ALTER TABLE candidates ADD COLUMN IF NOT EXISTS skills TEXT`,
		`ALTER TABLE candidates ADD COLUMN IF NOT EXISTS resume_url VARCHAR(500)`,
		`ALTER TABLE candidates ADD COLUMN IF NOT EXISTS linkedin_url VARCHAR(500)`,
		`ALTER TABLE candidates ADD COLUMN IF NOT EXISTS github_url VARCHAR(500)`,
		`ALTER TABLE candidates ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'active'`,
		`ALTER TABLE candidates ADD COLUMN IF NOT EXISTS interviewer_id UUID`,
		// Make email nullable if it's not already
		`ALTER TABLE candidates ALTER COLUMN email DROP NOT NULL`,
	}

	for _, query := range alterTableQueries {
		if _, err := db.Exec(query); err != nil {
			// Ignore errors for columns that don't exist or constraints that don't exist
			// This allows the function to work even if some changes have already been applied
		}
	}

	// Create unique constraint on email only for non-null values
	_, _ = db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS candidates_email_unique 
		ON candidates (email) 
		WHERE email IS NOT NULL
	`)

	return nil
}

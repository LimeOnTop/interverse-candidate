package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/LimeOnTop/interverse-candidate/internal/entity"
	"github.com/google/uuid"
)

type CandidateRepository struct {
	db *sql.DB
}

func NewCandidateRepository(db *sql.DB) *CandidateRepository {
	return &CandidateRepository{db: db}
}

func (r *CandidateRepository) Create(ctx context.Context, candidate entity.Candidate) (entity.Candidate, error) {
	candidate.ID = uuid.New().String()
	candidate.CreatedAt = time.Now()
	candidate.UpdatedAt = time.Now()

	email := nullableString(candidate.Email)

	query := `
		INSERT INTO candidates (id, name, email, phone, position, experience, skills, resume_url, linkedin_url, github_url, status, created_at, updated_at, interviewer_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.ExecContext(ctx, query,
		candidate.ID, candidate.Name, email, candidate.Phone,
		candidate.Position, candidate.Experience, candidate.Skills,
		candidate.ResumeURL, candidate.LinkedinURL, candidate.GithubURL,
		candidate.Status, candidate.CreatedAt, candidate.UpdatedAt, candidate.InterviewerID,
	)
	if err != nil {
		return entity.Candidate{}, fmt.Errorf("create candidate: %w", err)
	}

	return candidate, nil
}

func (r *CandidateRepository) GetByID(ctx context.Context, id string) (entity.Candidate, error) {
	query := `
		SELECT id, name, email, phone, position, experience, skills, resume_url, linkedin_url, github_url, status, created_at, updated_at, interviewer_id
		FROM candidates WHERE id = $1
	`

	candidate, err := scanCandidate(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Candidate{}, fmt.Errorf("get candidate: not found")
		}
		return entity.Candidate{}, fmt.Errorf("get candidate: %w", err)
	}

	return candidate, nil
}

func (r *CandidateRepository) GetByInterviewer(ctx context.Context, interviewerID string, limit, offset int) ([]entity.Candidate, error) {
	query := `
		SELECT id, name, email, phone, position, experience, skills, resume_url, linkedin_url, github_url, status, created_at, updated_at, interviewer_id
		FROM candidates
		WHERE interviewer_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, interviewerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get candidates: %w", err)
	}
	defer rows.Close()

	return scanCandidates(rows)
}

func (r *CandidateRepository) Update(ctx context.Context, candidate entity.Candidate) (entity.Candidate, error) {
	candidate.UpdatedAt = time.Now()
	email := nullableString(candidate.Email)

	query := `
		UPDATE candidates
		SET name = $1, email = $2, phone = $3, position = $4, experience = $5, skills = $6, resume_url = $7, linkedin_url = $8, github_url = $9, status = $10, updated_at = $11
		WHERE id = $12
	`

	result, err := r.db.ExecContext(ctx, query,
		candidate.Name, email, candidate.Phone, candidate.Position,
		candidate.Experience, candidate.Skills, candidate.ResumeURL,
		candidate.LinkedinURL, candidate.GithubURL, candidate.Status,
		candidate.UpdatedAt, candidate.ID,
	)
	if err != nil {
		return entity.Candidate{}, fmt.Errorf("update candidate: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return entity.Candidate{}, fmt.Errorf("update candidate: %w", err)
	}
	if rowsAffected == 0 {
		return entity.Candidate{}, fmt.Errorf("update candidate: not found")
	}

	return candidate, nil
}

func (r *CandidateRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM candidates WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete candidate: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete candidate: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("delete candidate: not found")
	}

	return nil
}

func (r *CandidateRepository) Search(ctx context.Context, queryText, interviewerID string, limit, offset int) ([]entity.Candidate, error) {
	query := `
		SELECT id, name, email, phone, position, experience, skills, resume_url, linkedin_url, github_url, status, created_at, updated_at, interviewer_id
		FROM candidates
		WHERE interviewer_id = $1 AND (
			LOWER(name) LIKE LOWER($2) OR
			LOWER(email) LIKE LOWER($2) OR
			LOWER(position) LIKE LOWER($2) OR
			LOWER(skills) LIKE LOWER($2)
		)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, interviewerID, "%"+queryText+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search candidates: %w", err)
	}
	defer rows.Close()

	return scanCandidates(rows)
}

func nullableString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}

func scanCandidate(row *sql.Row) (entity.Candidate, error) {
	var candidate entity.Candidate
	var email sql.NullString

	err := row.Scan(
		&candidate.ID, &candidate.Name, &email, &candidate.Phone,
		&candidate.Position, &candidate.Experience, &candidate.Skills,
		&candidate.ResumeURL, &candidate.LinkedinURL, &candidate.GithubURL,
		&candidate.Status, &candidate.CreatedAt, &candidate.UpdatedAt, &candidate.InterviewerID,
	)
	if err != nil {
		return entity.Candidate{}, err
	}

	if email.Valid {
		candidate.Email = email.String
	}

	return candidate, nil
}

func scanCandidates(rows *sql.Rows) ([]entity.Candidate, error) {
	var candidates []entity.Candidate

	for rows.Next() {
		var candidate entity.Candidate
		var email sql.NullString

		if err := rows.Scan(
			&candidate.ID, &candidate.Name, &email, &candidate.Phone,
			&candidate.Position, &candidate.Experience, &candidate.Skills,
			&candidate.ResumeURL, &candidate.LinkedinURL, &candidate.GithubURL,
			&candidate.Status, &candidate.CreatedAt, &candidate.UpdatedAt, &candidate.InterviewerID,
		); err != nil {
			return nil, fmt.Errorf("scan candidate: %w", err)
		}

		if email.Valid {
			candidate.Email = email.String
		}

		candidates = append(candidates, candidate)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan candidates: %w", err)
	}

	return candidates, nil
}

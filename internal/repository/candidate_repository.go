package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/LimeOnTop/interverse-candidate/internal/models"
)

type CandidateRepository struct {
	db *sql.DB
}

func NewCandidateRepository(db *sql.DB) *CandidateRepository {
	return &CandidateRepository{db: db}
}

func (r *CandidateRepository) CreateCandidate(candidate *models.Candidate) error {
	candidate.ID = uuid.New().String()
	candidate.CreatedAt = time.Now()
	candidate.UpdatedAt = time.Now()

	// Convert empty email to NULL to avoid unique constraint violations
	var email sql.NullString
	if candidate.Email != "" {
		email = sql.NullString{String: candidate.Email, Valid: true}
	} else {
		email = sql.NullString{Valid: false}
	}

	query := `
		INSERT INTO candidates (id, name, email, phone, position, experience, skills, resume_url, linkedin_url, github_url, status, created_at, updated_at, interviewer_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.Exec(query,
		candidate.ID, candidate.Name, email, candidate.Phone,
		candidate.Position, candidate.Experience, candidate.Skills,
		candidate.ResumeURL, candidate.LinkedinURL, candidate.GithubURL,
		candidate.Status, candidate.CreatedAt, candidate.UpdatedAt, candidate.InterviewerID)
	if err != nil {
		return fmt.Errorf("failed to create candidate: %w", err)
	}

	return nil
}

func (r *CandidateRepository) GetCandidateByID(id string) (*models.Candidate, error) {
	query := `
		SELECT id, name, email, phone, position, experience, skills, resume_url, linkedin_url, github_url, status, created_at, updated_at, interviewer_id
		FROM candidates WHERE id = $1
	`

	candidate := &models.Candidate{}
	var email sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&candidate.ID, &candidate.Name, &email, &candidate.Phone,
		&candidate.Position, &candidate.Experience, &candidate.Skills,
		&candidate.ResumeURL, &candidate.LinkedinURL, &candidate.GithubURL,
		&candidate.Status, &candidate.CreatedAt, &candidate.UpdatedAt, &candidate.InterviewerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("candidate not found")
		}
		return nil, fmt.Errorf("failed to get candidate: %w", err)
	}

	// Convert NULL email to empty string
	if email.Valid {
		candidate.Email = email.String
	} else {
		candidate.Email = ""
	}

	return candidate, nil
}

func (r *CandidateRepository) GetCandidatesByInterviewer(interviewerID string, limit, offset int) ([]*models.Candidate, error) {
	query := `
		SELECT id, name, email, phone, position, experience, skills, resume_url, linkedin_url, github_url, status, created_at, updated_at, interviewer_id
		FROM candidates 
		WHERE interviewer_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, interviewerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get candidates: %w", err)
	}
	defer rows.Close()

	var candidates []*models.Candidate
	for rows.Next() {
		candidate := &models.Candidate{}
		var email sql.NullString
		err := rows.Scan(
			&candidate.ID, &candidate.Name, &email, &candidate.Phone,
			&candidate.Position, &candidate.Experience, &candidate.Skills,
			&candidate.ResumeURL, &candidate.LinkedinURL, &candidate.GithubURL,
			&candidate.Status, &candidate.CreatedAt, &candidate.UpdatedAt, &candidate.InterviewerID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan candidate: %w", err)
		}
		// Convert NULL email to empty string
		if email.Valid {
			candidate.Email = email.String
		} else {
			candidate.Email = ""
		}
		candidates = append(candidates, candidate)
	}

	return candidates, nil
}

func (r *CandidateRepository) UpdateCandidate(candidate *models.Candidate) error {
	candidate.UpdatedAt = time.Now()

	// Convert empty email to NULL to avoid unique constraint violations
	var email sql.NullString
	if candidate.Email != "" {
		email = sql.NullString{String: candidate.Email, Valid: true}
	} else {
		email = sql.NullString{Valid: false}
	}

	query := `
		UPDATE candidates 
		SET name = $1, email = $2, phone = $3, position = $4, experience = $5, skills = $6, resume_url = $7, linkedin_url = $8, github_url = $9, status = $10, updated_at = $11
		WHERE id = $12
	`

	_, err := r.db.Exec(query,
		candidate.Name, email, candidate.Phone, candidate.Position,
		candidate.Experience, candidate.Skills, candidate.ResumeURL,
		candidate.LinkedinURL, candidate.GithubURL, candidate.Status,
		candidate.UpdatedAt, candidate.ID)
	if err != nil {
		return fmt.Errorf("failed to update candidate: %w", err)
	}

	return nil
}

func (r *CandidateRepository) DeleteCandidate(id string) error {
	query := `DELETE FROM candidates WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete candidate: %w", err)
	}

	return nil
}

func (r *CandidateRepository) SearchCandidates(query, interviewerID string, limit, offset int) ([]*models.Candidate, error) {
	searchQuery := `
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

	searchTerm := "%" + query + "%"
	rows, err := r.db.Query(searchQuery, interviewerID, searchTerm, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search candidates: %w", err)
	}
	defer rows.Close()

	var candidates []*models.Candidate
	for rows.Next() {
		candidate := &models.Candidate{}
		var email sql.NullString
		err := rows.Scan(
			&candidate.ID, &candidate.Name, &email, &candidate.Phone,
			&candidate.Position, &candidate.Experience, &candidate.Skills,
			&candidate.ResumeURL, &candidate.LinkedinURL, &candidate.GithubURL,
			&candidate.Status, &candidate.CreatedAt, &candidate.UpdatedAt, &candidate.InterviewerID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan candidate: %w", err)
		}
		// Convert NULL email to empty string
		if email.Valid {
			candidate.Email = email.String
		} else {
			candidate.Email = ""
		}
		candidates = append(candidates, candidate)
	}

	return candidates, nil
}

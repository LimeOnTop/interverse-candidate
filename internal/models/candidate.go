package models

import (
	"time"
)

type Candidate struct {
	ID            string    `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Email         string    `json:"email" db:"email"`
	Phone         string    `json:"phone" db:"phone"`
	Position      string    `json:"position" db:"position"`
	Experience    string    `json:"experience" db:"experience"`
	Skills        string    `json:"skills" db:"skills"`
	ResumeURL     string    `json:"resume_url" db:"resume_url"`
	LinkedinURL   string    `json:"linkedin_url" db:"linkedin_url"`
	GithubURL     string    `json:"github_url" db:"github_url"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	InterviewerID string    `json:"interviewer_id" db:"interviewer_id"`
}


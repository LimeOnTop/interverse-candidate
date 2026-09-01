package usecase

import "time"

type CandidateDTO struct {
	ID            string
	Name          string
	Email         string
	Phone         string
	Position      string
	Experience    int64
	Skills        string
	ResumeURL     string
	LinkedinURL   string
	GithubURL     string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	InterviewerID string
}

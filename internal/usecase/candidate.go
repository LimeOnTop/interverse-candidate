package usecase

import (
	"context"

	"github.com/LimeOnTop/interverse-candidate/internal/entity"
)

type Candidate interface {
	Create(ctx context.Context, candidate entity.Candidate) (CandidateDTO, error)
	GetByID(ctx context.Context, id string) (CandidateDTO, error)
	GetByInterviewer(ctx context.Context, interviewerID string, limit, offset int) ([]CandidateDTO, error)
	Update(ctx context.Context, candidate entity.Candidate) (CandidateDTO, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, query, interviewerID string, limit, offset int) ([]CandidateDTO, error)
}

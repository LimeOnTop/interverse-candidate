package usecase

import (
	"context"

	"github.com/LimeOnTop/interverse-candidate/internal/entity"
)

type CandidateRepository interface {
	Create(ctx context.Context, candidate entity.Candidate) (entity.Candidate, error)
	GetByID(ctx context.Context, id string) (entity.Candidate, error)
	GetByInterviewer(ctx context.Context, interviewerID string, limit, offset int) ([]entity.Candidate, error)
	Update(ctx context.Context, candidate entity.Candidate) (entity.Candidate, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, query, interviewerID string, limit, offset int) ([]entity.Candidate, error)
}

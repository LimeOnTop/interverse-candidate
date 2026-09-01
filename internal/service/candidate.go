package service

import (
	"context"
	"fmt"
	"time"

	"github.com/LimeOnTop/interverse-candidate/internal/entity"
	"github.com/LimeOnTop/interverse-candidate/internal/usecase"
)

const defaultStatus = "active"

type CandidateService struct {
	repository usecase.CandidateRepository
}

func NewCandidateService(repository usecase.CandidateRepository) *CandidateService {
	return &CandidateService{repository: repository}
}

var _ usecase.Candidate = (*CandidateService)(nil)

func (s *CandidateService) Create(ctx context.Context, candidate entity.Candidate) (usecase.CandidateDTO, error) {
	if candidate.Status == "" {
		candidate.Status = defaultStatus
	}
	candidate.CreatedAt = time.Now()
	candidate.UpdatedAt = time.Now()

	created, err := s.repository.Create(ctx, candidate)
	if err != nil {
		return usecase.CandidateDTO{}, fmt.Errorf("create candidate: %w", err)
	}

	return toDTO(created), nil
}

func (s *CandidateService) GetByID(ctx context.Context, id string) (usecase.CandidateDTO, error) {
	candidate, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return usecase.CandidateDTO{}, fmt.Errorf("get candidate: %w", err)
	}

	return toDTO(candidate), nil
}

func (s *CandidateService) GetByInterviewer(ctx context.Context, interviewerID string, limit, offset int) ([]usecase.CandidateDTO, error) {
	candidates, err := s.repository.GetByInterviewer(ctx, interviewerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get candidates: %w", err)
	}

	return toDTOs(candidates), nil
}

func (s *CandidateService) Update(ctx context.Context, candidate entity.Candidate) (usecase.CandidateDTO, error) {
	candidate.UpdatedAt = time.Now()

	updated, err := s.repository.Update(ctx, candidate)
	if err != nil {
		return usecase.CandidateDTO{}, fmt.Errorf("update candidate: %w", err)
	}

	return toDTO(updated), nil
}

func (s *CandidateService) Delete(ctx context.Context, id string) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete candidate: %w", err)
	}
	return nil
}

func (s *CandidateService) Search(ctx context.Context, query, interviewerID string, limit, offset int) ([]usecase.CandidateDTO, error) {
	candidates, err := s.repository.Search(ctx, query, interviewerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search candidates: %w", err)
	}

	return toDTOs(candidates), nil
}

func toDTO(candidate entity.Candidate) usecase.CandidateDTO {
	return usecase.CandidateDTO{
		ID:            candidate.ID,
		Name:          candidate.Name,
		Email:         candidate.Email,
		Phone:         candidate.Phone,
		Position:      candidate.Position,
		Experience:    candidate.Experience,
		Skills:        candidate.Skills,
		ResumeURL:     candidate.ResumeURL,
		LinkedinURL:   candidate.LinkedinURL,
		GithubURL:     candidate.GithubURL,
		Status:        candidate.Status,
		CreatedAt:     candidate.CreatedAt,
		UpdatedAt:     candidate.UpdatedAt,
		InterviewerID: candidate.InterviewerID,
	}
}

func toDTOs(candidates []entity.Candidate) []usecase.CandidateDTO {
	result := make([]usecase.CandidateDTO, 0, len(candidates))
	for _, candidate := range candidates {
		result = append(result, toDTO(candidate))
	}
	return result
}

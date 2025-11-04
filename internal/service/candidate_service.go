package service

import (
	"fmt"

	"github.com/inter-verse/candidate-service/internal/models"
	"github.com/inter-verse/candidate-service/internal/repository"
)

type CandidateService struct {
	candidateRepo  repository.CandidateRepository
	userServiceURL string
}

func NewCandidateService(candidateRepo repository.CandidateRepository, userServiceURL string) *CandidateService {
	return &CandidateService{
		candidateRepo:  candidateRepo,
		userServiceURL: userServiceURL,
	}
}

func (s *CandidateService) CreateCandidate(candidate *models.Candidate) (*models.Candidate, error) {
	if err := s.candidateRepo.CreateCandidate(candidate); err != nil {
		return nil, fmt.Errorf("failed to create candidate: %w", err)
	}

	return candidate, nil
}

func (s *CandidateService) GetCandidateByID(id string) (*models.Candidate, error) {
	candidate, err := s.candidateRepo.GetCandidateByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get candidate: %w", err)
	}

	return candidate, nil
}

func (s *CandidateService) GetCandidatesByInterviewer(interviewerID string, limit, offset int) ([]*models.Candidate, error) {
	candidates, err := s.candidateRepo.GetCandidatesByInterviewer(interviewerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get candidates: %w", err)
	}

	return candidates, nil
}

func (s *CandidateService) UpdateCandidate(candidate *models.Candidate) (*models.Candidate, error) {
	if err := s.candidateRepo.UpdateCandidate(candidate); err != nil {
		return nil, fmt.Errorf("failed to update candidate: %w", err)
	}

	return candidate, nil
}

func (s *CandidateService) DeleteCandidate(id string) error {
	return s.candidateRepo.DeleteCandidate(id)
}

func (s *CandidateService) SearchCandidates(query, interviewerID string, limit, offset int) ([]*models.Candidate, error) {
	candidates, err := s.candidateRepo.SearchCandidates(query, interviewerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search candidates: %w", err)
	}

	return candidates, nil
}

func (s *CandidateService) ValidateUser(userID string) error {
	// In a real implementation, this would call the User Service via gRPC
	// For now, we'll assume the user is valid
	return nil
}


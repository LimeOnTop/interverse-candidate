package handler

import (
	"context"
	"time"

	pb "github.com/LimeOnTop/interverse-contracts/candidate/gen"
	"github.com/LimeOnTop/interverse-candidate/internal/models"
	"github.com/LimeOnTop/interverse-candidate/internal/service"
)

type CandidateHandler struct {
	pb.UnimplementedCandidateServiceServer
	candidateService *service.CandidateService
}

func NewCandidateHandler(candidateService *service.CandidateService) *CandidateHandler {
	return &CandidateHandler{
		candidateService: candidateService,
	}
}

func (h *CandidateHandler) CreateCandidate(ctx context.Context, req *pb.CreateCandidateRequest) (*pb.CreateCandidateResponse, error) {
	candidate := &models.Candidate{
		Name:          req.Name,
		Email:         req.Email,
		Phone:         req.Phone,
		Position:      req.Position,
		Experience:    req.Experience,
		Skills:        req.Skills,
		ResumeURL:     req.ResumeUrl,
		LinkedinURL:   req.LinkedinUrl,
		GithubURL:     req.GithubUrl,
		Status:        "active",
		InterviewerID: req.InterviewerId,
	}

	createdCandidate, err := h.candidateService.CreateCandidate(candidate)
	if err != nil {
		return &pb.CreateCandidateResponse{
			Response: &pb.Response{
				Success: false,
				Error:   err.Error(),
			},
		}, nil
	}

	return &pb.CreateCandidateResponse{
		Response: &pb.Response{
			Success: true,
			Message: "Candidate created successfully",
		},
		Candidate: &pb.Candidate{
			Id:            createdCandidate.ID,
			Name:          createdCandidate.Name,
			Email:         createdCandidate.Email,
			Phone:         createdCandidate.Phone,
			Position:      createdCandidate.Position,
			Experience:    createdCandidate.Experience,
			Skills:        createdCandidate.Skills,
			ResumeUrl:     createdCandidate.ResumeURL,
			LinkedinUrl:   createdCandidate.LinkedinURL,
			GithubUrl:     createdCandidate.GithubURL,
			Status:        createdCandidate.Status,
			CreatedAt:     createdCandidate.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     createdCandidate.UpdatedAt.Format(time.RFC3339),
			InterviewerId: createdCandidate.InterviewerID,
		},
	}, nil
}

func (h *CandidateHandler) GetCandidate(ctx context.Context, req *pb.GetCandidateRequest) (*pb.GetCandidateResponse, error) {
	candidate, err := h.candidateService.GetCandidateByID(req.CandidateId)
	if err != nil {
		return &pb.GetCandidateResponse{
			Response: &pb.Response{
				Success: false,
				Error:   err.Error(),
			},
		}, nil
	}

	return &pb.GetCandidateResponse{
		Response: &pb.Response{
			Success: true,
		},
		Candidate: &pb.Candidate{
			Id:            candidate.ID,
			Name:          candidate.Name,
			Email:         candidate.Email,
			Phone:         candidate.Phone,
			Position:      candidate.Position,
			Experience:    candidate.Experience,
			Skills:        candidate.Skills,
			ResumeUrl:     candidate.ResumeURL,
			LinkedinUrl:   candidate.LinkedinURL,
			GithubUrl:     candidate.GithubURL,
			Status:        candidate.Status,
			CreatedAt:     candidate.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     candidate.UpdatedAt.Format(time.RFC3339),
			InterviewerId: candidate.InterviewerID,
		},
	}, nil
}

func (h *CandidateHandler) GetCandidates(ctx context.Context, req *pb.GetCandidatesRequest) (*pb.GetCandidatesResponse, error) {
	limit := 10
	offset := 0

	if req.Pagination != nil {
		limit = int(req.Pagination.Limit)
		offset = int(req.Pagination.Page-1) * int(req.Pagination.Limit)
	}

	candidates, err := h.candidateService.GetCandidatesByInterviewer(req.InterviewerId, limit, offset)
	if err != nil {
		return &pb.GetCandidatesResponse{
			Response: &pb.Response{
				Success: false,
				Error:   err.Error(),
			},
		}, nil
	}

	var pbCandidates []*pb.Candidate
	for _, candidate := range candidates {
		pbCandidates = append(pbCandidates, &pb.Candidate{
			Id:            candidate.ID,
			Name:          candidate.Name,
			Email:         candidate.Email,
			Phone:         candidate.Phone,
			Position:      candidate.Position,
			Experience:    candidate.Experience,
			Skills:        candidate.Skills,
			ResumeUrl:     candidate.ResumeURL,
			LinkedinUrl:   candidate.LinkedinURL,
			GithubUrl:     candidate.GithubURL,
			Status:        candidate.Status,
			CreatedAt:     candidate.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     candidate.UpdatedAt.Format(time.RFC3339),
			InterviewerId: candidate.InterviewerID,
		})
	}

	return &pb.GetCandidatesResponse{
		Response: &pb.Response{
			Success: true,
		},
		Candidates: pbCandidates,
		Pagination: &pb.Pagination{
			Page:  req.Pagination.Page,
			Limit: req.Pagination.Limit,
			Total: int32(len(pbCandidates)),
		},
	}, nil
}

func (h *CandidateHandler) UpdateCandidate(ctx context.Context, req *pb.UpdateCandidateRequest) (*pb.UpdateCandidateResponse, error) {
	candidate := &models.Candidate{
		ID:          req.CandidateId,
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Position:    req.Position,
		Experience:  req.Experience,
		Skills:      req.Skills,
		ResumeURL:   req.ResumeUrl,
		LinkedinURL: req.LinkedinUrl,
		GithubURL:   req.GithubUrl,
		Status:      req.Status,
	}

	updatedCandidate, err := h.candidateService.UpdateCandidate(candidate)
	if err != nil {
		return &pb.UpdateCandidateResponse{
			Response: &pb.Response{
				Success: false,
				Error:   err.Error(),
			},
		}, nil
	}

	return &pb.UpdateCandidateResponse{
		Response: &pb.Response{
			Success: true,
			Message: "Candidate updated successfully",
		},
		Candidate: &pb.Candidate{
			Id:            updatedCandidate.ID,
			Name:          updatedCandidate.Name,
			Email:         updatedCandidate.Email,
			Phone:         updatedCandidate.Phone,
			Position:      updatedCandidate.Position,
			Experience:    updatedCandidate.Experience,
			Skills:        updatedCandidate.Skills,
			ResumeUrl:     updatedCandidate.ResumeURL,
			LinkedinUrl:   updatedCandidate.LinkedinURL,
			GithubUrl:     updatedCandidate.GithubURL,
			Status:        updatedCandidate.Status,
			CreatedAt:     updatedCandidate.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     updatedCandidate.UpdatedAt.Format(time.RFC3339),
			InterviewerId: updatedCandidate.InterviewerID,
		},
	}, nil
}

func (h *CandidateHandler) DeleteCandidate(ctx context.Context, req *pb.DeleteCandidateRequest) (*pb.Response, error) {
	err := h.candidateService.DeleteCandidate(req.CandidateId)
	if err != nil {
		return &pb.Response{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.Response{
		Success: true,
		Message: "Candidate deleted successfully",
	}, nil
}

func (h *CandidateHandler) SearchCandidates(ctx context.Context, req *pb.SearchCandidatesRequest) (*pb.SearchCandidatesResponse, error) {
	limit := 10
	offset := 0

	if req.Pagination != nil {
		limit = int(req.Pagination.Limit)
		offset = int(req.Pagination.Page-1) * int(req.Pagination.Limit)
	}

	candidates, err := h.candidateService.SearchCandidates(req.Query, req.InterviewerId, limit, offset)
	if err != nil {
		return &pb.SearchCandidatesResponse{
			Response: &pb.Response{
				Success: false,
				Error:   err.Error(),
			},
		}, nil
	}

	var pbCandidates []*pb.Candidate
	for _, candidate := range candidates {
		pbCandidates = append(pbCandidates, &pb.Candidate{
			Id:            candidate.ID,
			Name:          candidate.Name,
			Email:         candidate.Email,
			Phone:         candidate.Phone,
			Position:      candidate.Position,
			Experience:    candidate.Experience,
			Skills:        candidate.Skills,
			ResumeUrl:     candidate.ResumeURL,
			LinkedinUrl:   candidate.LinkedinURL,
			GithubUrl:     candidate.GithubURL,
			Status:        candidate.Status,
			CreatedAt:     candidate.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     candidate.UpdatedAt.Format(time.RFC3339),
			InterviewerId: candidate.InterviewerID,
		})
	}

	return &pb.SearchCandidatesResponse{
		Response: &pb.Response{
			Success: true,
		},
		Candidates: pbCandidates,
		Pagination: &pb.Pagination{
			Page:  req.Pagination.Page,
			Limit: req.Pagination.Limit,
			Total: int32(len(pbCandidates)),
		},
	}, nil
}

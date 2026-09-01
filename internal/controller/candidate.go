package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/LimeOnTop/interverse-candidate/internal/entity"
	"github.com/LimeOnTop/interverse-candidate/internal/usecase"
	pb "github.com/LimeOnTop/interverse-contracts/candidate/gen"
)

type CandidateController struct {
	pb.UnimplementedCandidateServiceServer
	candidate usecase.Candidate
}

func NewCandidateController(candidate usecase.Candidate) *CandidateController {
	return &CandidateController{candidate: candidate}
}

func (c *CandidateController) CreateCandidate(ctx context.Context, req *pb.CreateCandidateRequest) (*pb.CreateCandidateResponse, error) {
	created, err := c.candidate.Create(ctx, entity.Candidate{
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
	})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("context canceled: %w", err)
		}
		return &pb.CreateCandidateResponse{
			Response: &pb.Response{
				Success: false,
				Message: "Create candidate failed",
				Error:   err.Error(),
			},
		}, nil
	}

	return &pb.CreateCandidateResponse{
		Response: &pb.Response{
			Success: true,
			Message: "Candidate created successfully",
		},
		Candidate: toProtoCandidate(created),
	}, nil
}

func (c *CandidateController) GetCandidate(ctx context.Context, req *pb.GetCandidateRequest) (*pb.GetCandidateResponse, error) {
	candidate, err := c.candidate.GetByID(ctx, req.CandidateId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("context canceled: %w", err)
		}
		return &pb.GetCandidateResponse{
			Response: &pb.Response{
				Success: false,
				Message: "Get candidate failed",
				Error:   err.Error(),
			},
		}, nil
	}

	return &pb.GetCandidateResponse{
		Response: &pb.Response{
			Success: true,
			Message: "Get candidate completed",
		},
		Candidate: toProtoCandidate(candidate),
	}, nil
}

func (c *CandidateController) GetCandidates(ctx context.Context, req *pb.GetCandidatesRequest) (*pb.GetCandidatesResponse, error) {
	limit := 10
	offset := 0

	if req.Pagination != nil {
		limit = int(req.Pagination.Limit)
		offset = int(req.Pagination.Page-1) * int(req.Pagination.Limit)
	}

	candidates, err := c.candidate.GetByInterviewer(ctx, req.InterviewerId, limit, offset)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("context canceled: %w", err)
		}
		return &pb.GetCandidatesResponse{
			Response: &pb.Response{
				Success: false,
				Message: "Get candidates failed",
				Error:   err.Error(),
			},
		}, nil
	}

	return &pb.GetCandidatesResponse{
		Response: &pb.Response{
			Success: true,
			Message: "Get candidates completed",
		},
		Candidates: toProtoCandidates(candidates),
		Pagination: &pb.Pagination{
			Page:  req.Pagination.GetPage(),
			Limit: req.Pagination.GetLimit(),
			Total: int32(len(candidates)),
		},
	}, nil
}

func (c *CandidateController) UpdateCandidate(ctx context.Context, req *pb.UpdateCandidateRequest) (*pb.UpdateCandidateResponse, error) {
	updated, err := c.candidate.Update(ctx, entity.Candidate{
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
	})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("context canceled: %w", err)
		}
		return &pb.UpdateCandidateResponse{
			Response: &pb.Response{
				Success: false,
				Message: "Update candidate failed",
				Error:   err.Error(),
			},
		}, nil
	}

	return &pb.UpdateCandidateResponse{
		Response: &pb.Response{
			Success: true,
			Message: "Candidate updated successfully",
		},
		Candidate: toProtoCandidate(updated),
	}, nil
}

func (c *CandidateController) DeleteCandidate(ctx context.Context, req *pb.DeleteCandidateRequest) (*pb.Response, error) {
	err := c.candidate.Delete(ctx, req.CandidateId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("context canceled: %w", err)
		}
		return &pb.Response{
			Success: false,
			Message: "Delete candidate failed",
			Error:   err.Error(),
		}, nil
	}

	return &pb.Response{
		Success: true,
		Message: "Candidate deleted successfully",
	}, nil
}

func (c *CandidateController) SearchCandidates(ctx context.Context, req *pb.SearchCandidatesRequest) (*pb.SearchCandidatesResponse, error) {
	limit := 10
	offset := 0

	if req.Pagination != nil {
		limit = int(req.Pagination.Limit)
		offset = int(req.Pagination.Page-1) * int(req.Pagination.Limit)
	}

	candidates, err := c.candidate.Search(ctx, req.Query, req.InterviewerId, limit, offset)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("context canceled: %w", err)
		}
		return &pb.SearchCandidatesResponse{
			Response: &pb.Response{
				Success: false,
				Message: "Search candidates failed",
				Error:   err.Error(),
			},
		}, nil
	}

	return &pb.SearchCandidatesResponse{
		Response: &pb.Response{
			Success: true,
			Message: "Search candidates completed",
		},
		Candidates: toProtoCandidates(candidates),
		Pagination: &pb.Pagination{
			Page:  req.Pagination.GetPage(),
			Limit: req.Pagination.GetLimit(),
			Total: int32(len(candidates)),
		},
	}, nil
}

func toProtoCandidate(candidate usecase.CandidateDTO) *pb.Candidate {
	return &pb.Candidate{
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
	}
}

func toProtoCandidates(candidates []usecase.CandidateDTO) []*pb.Candidate {
	result := make([]*pb.Candidate, 0, len(candidates))
	for _, candidate := range candidates {
		result = append(result, toProtoCandidate(candidate))
	}
	return result
}

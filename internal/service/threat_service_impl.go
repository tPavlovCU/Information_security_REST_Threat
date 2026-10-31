package service

import (
	"context"

	"threat-intelligence/internal/domain"
	"threat-intelligence/internal/repository"
)

type threatService struct {
	repo repository.ThreatRepository
}

func NewThreatService(repo repository.ThreatRepository) ThreatService {
	return &threatService{
		repo: repo,
	}
}

func (s *threatService) GetByMD5(ctx context.Context, targetMD5 string) (*domain.Threat, error) {
	return s.repo.GetByMD5(ctx, targetMD5)
}

func (s *threatService) Create(ctx context.Context, threat *domain.Threat) error {
	return nil
}

func (s *threatService) Update(ctx context.Context, threat *domain.Threat) error {
	return nil
}

var _ ThreatService = (*threatService)(nil)

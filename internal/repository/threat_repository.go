package repository

import (
	"context"
	"threat-intelligence/internal/domain"
)

type ThreatRepository interface {
	Create(ctx context.Context, threat *domain.Threat) error
	GetByMD5(ctx context.Context, targetMD5 string) (*domain.Threat, error)
	Update(ctx context.Context, threat *domain.Threat) error
}

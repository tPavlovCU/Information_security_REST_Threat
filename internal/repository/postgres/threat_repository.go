package postgres

import (
	"context"
	"errors"

	"threat-intelligence/internal/domain"
	"threat-intelligence/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ThreatRepository struct {
	pool *pgxpool.Pool
}

func NewThreatRepository(pool *pgxpool.Pool) *ThreatRepository {
	return &ThreatRepository{pool}
}

func (r *ThreatRepository) GetByMD5(ctx context.Context, targetMD5 string) (*domain.Threat, error) {
	query := `
	SELECT 
	id, 
	md5, 
	sha256,
	name,
	family, 
	type, 
	severity, 
	confidence, 
	description, 
	tags, 
	source 
	FROM threats
	WHERE md5 = $1
	`
	var (
		id          int64
		md5         string
		sha256      *string
		name        string
		family      *string
		resType     string
		severity    string
		confidence  *int16
		description *string
		tags        []string
		source      *string
	)

	err := r.pool.QueryRow(ctx, query, targetMD5).Scan(
		&id,
		&md5,
		&sha256,
		&name,
		&family,
		&resType,
		&severity,
		&confidence,
		&description,
		&tags,
		&source,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrThreatNotFound
		}

		return nil, err
	}

	threat := domain.Threat{}
	threat.ID = id
	threat.MD5 = md5
	threat.SHA256 = sha256
	threat.Name = name
	threat.Family = family
	threat.Type = resType
	threat.Severity = severity
	threat.Confidence = confidence
	threat.Description = description
	threat.Tags = tags
	threat.Source = source

	return &threat, nil
}

var _ repository.ThreatRepository = (*ThreatRepository)(nil)

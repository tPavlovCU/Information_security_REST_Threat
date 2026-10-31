package postgres

import (
	"context"
	"errors"

	"threat-intelligence/internal/domain"
	"threat-intelligence/internal/repository"

	"github.com/jackc/pgconn"
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

func (r *ThreatRepository) Create(ctx context.Context, threat *domain.Threat) error {
	query := `
	INSERT INTO threats (
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
	) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.pool.Exec(ctx, query,
		threat.MD5,
		threat.SHA256,
		threat.Name,
		threat.Family,
		threat.Type,
		threat.Severity,
		threat.Confidence,
		threat.Description,
		threat.Tags,
		threat.Source)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return domain.ErrThreatAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (r *ThreatRepository) Update(ctx context.Context, threat *domain.Threat) error {
	query := `
	UPDATE threats
	SET
		sha256 = $1,
		name = $2,
		family = $3,
		type = $4,
		severity = $5,
		confidence = $6,
		description = $7,
		tags = $8,
		source = $9,
		updated_at = NOW()
	WHERE md5 = $10
	`

	tag, err := r.pool.Exec(ctx, query,
		threat.SHA256,
		threat.Name,
		threat.Family,
		threat.Type,
		threat.Severity,
		threat.Confidence,
		threat.Description,
		threat.Tags,
		threat.Source,
		threat.MD5)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrThreatNotFound
	}
	return nil
}

var _ repository.ThreatRepository = (*ThreatRepository)(nil)

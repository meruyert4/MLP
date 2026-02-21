package audio

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, audio *Audio) error
	GetByID(ctx context.Context, id uuid.UUID) (*Audio, error)
	GetByLectureID(ctx context.Context, lectureID uuid.UUID) ([]Audio, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, audio *Audio) error {
	query := `
		INSERT INTO audios (id, lecture_id, url, language, voice, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, audio.ID, audio.LectureID, audio.URL, audio.Language, audio.Voice, audio.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create audio: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Audio, error) {
	query := `
		SELECT id, lecture_id, url, language, voice, created_at
		FROM audios
		WHERE id = $1
	`
	var audio Audio
	err := r.db.QueryRow(ctx, query, id).Scan(
		&audio.ID,
		&audio.LectureID,
		&audio.URL,
		&audio.Language,
		&audio.Voice,
		&audio.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get audio: %w", err)
	}
	return &audio, nil
}

func (r *repository) GetByLectureID(ctx context.Context, lectureID uuid.UUID) ([]Audio, error) {
	query := `
		SELECT id, lecture_id, url, language, voice, created_at
		FROM audios
		WHERE lecture_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, lectureID)
	if err != nil {
		return nil, fmt.Errorf("failed to get audios: %w", err)
	}
	defer rows.Close()

	var audios []Audio
	for rows.Next() {
		var audio Audio
		if err := rows.Scan(&audio.ID, &audio.LectureID, &audio.URL, &audio.Language, &audio.Voice, &audio.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan audio: %w", err)
		}
		audios = append(audios, audio)
	}

	return audios, nil
}

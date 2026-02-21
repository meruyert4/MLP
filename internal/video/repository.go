package video

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, video *Video) error
	GetByID(ctx context.Context, id uuid.UUID) (*Video, error)
	GetByAudioID(ctx context.Context, audioID uuid.UUID) ([]Video, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateURL(ctx context.Context, id uuid.UUID, url string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, video *Video) error {
	query := `
		INSERT INTO videos (id, audio_id, url, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, video.ID, video.AudioID, video.URL, video.Status, video.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create video: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Video, error) {
	query := `
		SELECT id, audio_id, url, status, created_at
		FROM videos
		WHERE id = $1
	`
	var video Video
	err := r.db.QueryRow(ctx, query, id).Scan(
		&video.ID,
		&video.AudioID,
		&video.URL,
		&video.Status,
		&video.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get video: %w", err)
	}
	return &video, nil
}

func (r *repository) GetByAudioID(ctx context.Context, audioID uuid.UUID) ([]Video, error) {
	query := `
		SELECT id, audio_id, url, status, created_at
		FROM videos
		WHERE audio_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, audioID)
	if err != nil {
		return nil, fmt.Errorf("failed to get videos: %w", err)
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		var video Video
		if err := rows.Scan(&video.ID, &video.AudioID, &video.URL, &video.Status, &video.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan video: %w", err)
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (r *repository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE videos SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

func (r *repository) UpdateURL(ctx context.Context, id uuid.UUID, url string) error {
	query := `UPDATE videos SET url = $1, status = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, url, "completed", id)
	return err
}

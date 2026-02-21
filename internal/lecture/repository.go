package lecture

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, lecture *Lecture) error
	GetByID(ctx context.Context, id uuid.UUID) (*Lecture, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Lecture, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, lecture *Lecture) error {
	query := `
		INSERT INTO lectures (id, user_id, topic, content, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, lecture.ID, lecture.UserID, lecture.Topic, lecture.Content, lecture.Status, lecture.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create lecture: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Lecture, error) {
	query := `
		SELECT id, user_id, topic, content, status, created_at
		FROM lectures
		WHERE id = $1
	`
	var lecture Lecture
	err := r.db.QueryRow(ctx, query, id).Scan(
		&lecture.ID,
		&lecture.UserID,
		&lecture.Topic,
		&lecture.Content,
		&lecture.Status,
		&lecture.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get lecture: %w", err)
	}
	return &lecture, nil
}

func (r *repository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]Lecture, error) {
	query := `
		SELECT id, user_id, topic, content, status, created_at
		FROM lectures
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lectures: %w", err)
	}
	defer rows.Close()

	var lectures []Lecture
	for rows.Next() {
		var lecture Lecture
		if err := rows.Scan(&lecture.ID, &lecture.UserID, &lecture.Topic, &lecture.Content, &lecture.Status, &lecture.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan lecture: %w", err)
		}
		lectures = append(lectures, lecture)
	}

	return lectures, nil
}

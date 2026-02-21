package lecture

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, req CreateLectureRequest) (*Lecture, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Lecture, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Lecture, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, req CreateLectureRequest) (*Lecture, error) {
	lecture := &Lecture{
		ID:        uuid.New(),
		UserID:    userID,
		Topic:     req.Topic,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, lecture); err != nil {
		return nil, err
	}

	return lecture, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Lecture, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByUserID(ctx context.Context, userID uuid.UUID) ([]Lecture, error) {
	return s.repo.GetByUserID(ctx, userID)
}

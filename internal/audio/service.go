package audio

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, lectureID uuid.UUID, url string) (*Audio, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Audio, error)
	GetByLectureID(ctx context.Context, lectureID uuid.UUID) ([]Audio, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, lectureID uuid.UUID, url string) (*Audio, error) {
	audio := &Audio{
		ID:        uuid.New(),
		LectureID: lectureID,
		URL:       url,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, audio); err != nil {
		return nil, err
	}

	return audio, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Audio, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByLectureID(ctx context.Context, lectureID uuid.UUID) ([]Audio, error) {
	return s.repo.GetByLectureID(ctx, lectureID)
}

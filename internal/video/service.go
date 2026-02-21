package video

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, audioID uuid.UUID, url string) (*Video, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Video, error)
	GetByAudioID(ctx context.Context, audioID uuid.UUID) ([]Video, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, audioID uuid.UUID, url string) (*Video, error) {
	video := &Video{
		ID:        uuid.New(),
		AudioID:   audioID,
		URL:       url,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, video); err != nil {
		return nil, err
	}

	return video, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Video, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByAudioID(ctx context.Context, audioID uuid.UUID) ([]Video, error) {
	return s.repo.GetByAudioID(ctx, audioID)
}

package lecture

import (
	"context"
	"time"

	"mlp/pkg/gemini"

	"github.com/google/uuid"
)

type Service interface {
	GenerateLecture(ctx context.Context, userID uuid.UUID, req CreateLectureRequest) (*GenerateLectureResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Lecture, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Lecture, error)
}

type service struct {
	repo          Repository
	geminiClient  *gemini.Client
}

func NewService(repo Repository, geminiClient *gemini.Client) Service {
	return &service{
		repo:         repo,
		geminiClient: geminiClient,
	}
}

func (s *service) GenerateLecture(ctx context.Context, userID uuid.UUID, req CreateLectureRequest) (*GenerateLectureResponse, error) {
	content, err := s.geminiClient.GenerateLecture(ctx, req.Topic)
	if err != nil {
		return nil, err
	}

	lecture := &Lecture{
		ID:        uuid.New(),
		UserID:    userID,
		Topic:     req.Topic,
		Content:   content,
		Status:    "completed",
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, lecture); err != nil {
		return nil, err
	}

	return &GenerateLectureResponse{
		ID:      lecture.ID,
		Topic:   lecture.Topic,
		Content: lecture.Content,
		Status:  lecture.Status,
	}, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Lecture, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByUserID(ctx context.Context, userID uuid.UUID) ([]Lecture, error) {
	return s.repo.GetByUserID(ctx, userID)
}

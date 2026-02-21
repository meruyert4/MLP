package audio

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"mlp/internal/lecture"
	"mlp/internal/storage"
	"mlp/pkg/voicerss"

	"github.com/google/uuid"
)

type Service interface {
	CreateFromLecture(ctx context.Context, req CreateAudioRequest) (*Audio, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Audio, error)
	GetByLectureID(ctx context.Context, lectureID uuid.UUID) ([]Audio, error)
}

type service struct {
	repo            Repository
	lectureRepo     lecture.Repository
	voiceRSSClient  *voicerss.Client
	minioClient     *storage.MinIOClient
	audioBucket     string
}

func NewService(repo Repository, lectureRepo lecture.Repository, voiceRSSClient *voicerss.Client, minioClient *storage.MinIOClient, audioBucket string) Service {
	return &service{
		repo:           repo,
		lectureRepo:    lectureRepo,
		voiceRSSClient: voiceRSSClient,
		minioClient:    minioClient,
		audioBucket:    audioBucket,
	}
}

func (s *service) CreateFromLecture(ctx context.Context, req CreateAudioRequest) (*Audio, error) {
	lect, err := s.lectureRepo.GetByID(ctx, req.LectureID)
	if err != nil {
		return nil, fmt.Errorf("lecture not found: %w", err)
	}

	if req.Language == "" {
		req.Language = "en-us"
	}

	audioData, err := s.voiceRSSClient.TextToSpeech(ctx, voicerss.TTSRequest{
		Text:     lect.Content,
		Language: req.Language,
		Voice:    req.Voice,
		Rate:     req.Rate,
		Codec:    "MP3",
		Format:   "16khz_16bit_stereo",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate audio: %w", err)
	}

	audioID := uuid.New()
	objectName := fmt.Sprintf("%s/%s.mp3", req.LectureID.String(), audioID.String())

	url, err := s.minioClient.Upload(
		ctx,
		s.audioBucket,
		objectName,
		bytes.NewReader(audioData),
		int64(len(audioData)),
		"audio/mpeg",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload audio: %w", err)
	}

	audio := &Audio{
		ID:        audioID,
		LectureID: req.LectureID,
		URL:       url,
		Language:  req.Language,
		Voice:     req.Voice,
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

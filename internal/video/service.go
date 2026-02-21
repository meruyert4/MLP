package video

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"mlp/internal/audio"
	"mlp/internal/storage"
	"mlp/pkg/lipsync"

	"github.com/google/uuid"
)

type Service interface {
	CreateFromAudio(ctx context.Context, req CreateVideoRequest, avatarFile io.Reader, avatarFilename string) (*Video, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Video, error)
	GetByAudioID(ctx context.Context, audioID uuid.UUID) ([]Video, error)
}

type service struct {
	repo           Repository
	audioRepo      audio.Repository
	minioClient    *storage.MinIOClient
	lipsyncClient  *lipsync.Client
	videoBucket    string
	avatarBucket   string
	minioEndpoint  string
}

func NewService(repo Repository, audioRepo audio.Repository, minioClient *storage.MinIOClient, lipsyncClient *lipsync.Client, videoBucket, avatarBucket, minioEndpoint string) Service {
	return &service{
		repo:          repo,
		audioRepo:     audioRepo,
		minioClient:   minioClient,
		lipsyncClient: lipsyncClient,
		videoBucket:   videoBucket,
		avatarBucket:  avatarBucket,
		minioEndpoint: minioEndpoint,
	}
}

func (s *service) CreateFromAudio(ctx context.Context, req CreateVideoRequest, avatarFile io.Reader, avatarFilename string) (*Video, error) {
	aud, err := s.audioRepo.GetByID(ctx, req.AudioID)
	if err != nil {
		return nil, fmt.Errorf("audio not found: %w", err)
	}

	avatarData, err := io.ReadAll(avatarFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read avatar file: %w", err)
	}

	avatarID := uuid.New()
	avatarObjectName := fmt.Sprintf("%s/%s", req.AudioID.String(), avatarID.String()+getFileExtension(avatarFilename))

	avatarURL, err := s.minioClient.Upload(
		ctx,
		s.avatarBucket,
		avatarObjectName,
		bytes.NewReader(avatarData),
		int64(len(avatarData)),
		"image/jpeg",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload avatar: %w", err)
	}

	videoID := uuid.New()
	video := &Video{
		ID:        videoID,
		AudioID:   req.AudioID,
		URL:       "",
		Status:    "processing",
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, video); err != nil {
		return nil, err
	}

	go s.processLipsync(context.Background(), videoID, avatarURL, aud.URL)

	return video, nil
}

func (s *service) processLipsync(ctx context.Context, videoID uuid.UUID, avatarURL, audioURL string) {
	avatarFullURL := fmt.Sprintf("http://%s%s", s.minioEndpoint, avatarURL)
	audioFullURL := fmt.Sprintf("http://%s%s", s.minioEndpoint, audioURL)

	syncResp, err := s.lipsyncClient.CreateLipsyncJob(ctx, avatarFullURL, audioFullURL)
	if err != nil {
		s.repo.UpdateStatus(ctx, videoID, "failed")
		return
	}

	for i := 0; i < 60; i++ {
		time.Sleep(5 * time.Second)

		status, err := s.lipsyncClient.GetJobStatus(ctx, syncResp.JobID)
		if err != nil {
			continue
		}

		if status.Status == "completed" && status.OutputURL != "" {
			s.repo.UpdateURL(ctx, videoID, status.OutputURL)
			return
		}

		if status.Status == "failed" {
			s.repo.UpdateStatus(ctx, videoID, "failed")
			return
		}
	}

	s.repo.UpdateStatus(ctx, videoID, "timeout")
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Video, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByAudioID(ctx context.Context, audioID uuid.UUID) ([]Video, error) {
	return s.repo.GetByAudioID(ctx, audioID)
}

func getFileExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}

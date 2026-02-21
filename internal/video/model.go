package video

import (
	"time"

	"github.com/google/uuid"
)

type Video struct {
	ID        uuid.UUID `json:"id"`
	AudioID   uuid.UUID `json:"audio_id"`
	URL       string    `json:"url"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateVideoRequest struct {
	AudioID uuid.UUID `json:"audio_id"`
}

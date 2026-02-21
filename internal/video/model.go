package video

import (
	"time"

	"github.com/google/uuid"
)

type Video struct {
	ID        uuid.UUID `json:"id"`
	AudioID   uuid.UUID `json:"audio_id"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

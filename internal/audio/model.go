package audio

import (
	"time"

	"github.com/google/uuid"
)

type Audio struct {
	ID        uuid.UUID `json:"id"`
	LectureID uuid.UUID `json:"lecture_id"`
	URL       string    `json:"url"`
	Language  string    `json:"language"`
	Voice     string    `json:"voice"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateAudioRequest struct {
	LectureID uuid.UUID `json:"lecture_id"`
	Language  string    `json:"language"`
	Voice     string    `json:"voice"`
	Rate      int       `json:"rate"`
}

package audio

import (
	"time"

	"github.com/google/uuid"
)

type Audio struct {
	ID        uuid.UUID `json:"id"`
	LectureID uuid.UUID `json:"lecture_id"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

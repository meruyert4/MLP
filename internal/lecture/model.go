package lecture

import (
	"time"

	"github.com/google/uuid"
)

type Lecture struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Topic     string    `json:"topic"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateLectureRequest struct {
	Topic   string `json:"topic"`
	Content string `json:"content"`
}

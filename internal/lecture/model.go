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
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateLectureRequest struct {
	Topic string `json:"topic"`
}

type GenerateLectureResponse struct {
	ID      uuid.UUID `json:"id"`
	Topic   string    `json:"topic"`
	Content string    `json:"content"`
	Status  string    `json:"status"`
}

// TestQuestion is a single question with 4 variants and the correct answer.
type TestQuestion struct {
	Question       string   `json:"question"`
	Variants       []string `json:"variants"`
	CorrectVariant string   `json:"correct_variant"`
}

// GenerateTestResponse is the response for generated test (10 questions).
type GenerateTestResponse struct {
	Questions []TestQuestion `json:"questions"`
}

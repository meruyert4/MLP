package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

type cohereRequest struct {
	Message string `json:"message"`
	Model   string `json:"model"`
	Stream  bool   `json:"stream"`
}

type cohereResponse struct {
	Text string `json:"text"`
}

func NewClient(apiKey, model string) *Client {
	return &Client{
		apiKey:  apiKey,
		model:   model,
		baseURL: "https://api.cohere.com/v1/chat",
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Client) GenerateLecture(ctx context.Context, topic string) (string, error) {
	prompt := fmt.Sprintf(`Create a comprehensive educational lecture on the topic: "%s". 

The lecture should:
1. Start with an engaging introduction
2. Cover key concepts and principles
3. Include relevant examples and explanations
4. Be structured in a logical flow
5. End with a conclusion and key takeaways

Write in a clear, educational tone suitable for text-to-speech conversion.
Aim for approximately 500-800 words.`, topic)

	reqBody := cohereRequest{
		Message: prompt,
		Model:   c.model,
		Stream:  false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var genResp cohereResponse
	if err := json.Unmarshal(body, &genResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if genResp.Text == "" {
		return "", fmt.Errorf("no content generated")
	}

	return genResp.Text, nil
}

// TestQuestion represents a single question with 4 variants and the correct answer.
type TestQuestion struct {
	Question       string   `json:"question"`
	Variants       []string `json:"variants"`        // 4 options
	CorrectVariant string   `json:"correct_variant"` // the correct answer text
}

// Test holds 10 questions generated from a lecture.
type Test struct {
	Questions []TestQuestion `json:"questions"`
}

// GenerateTest creates a test from lecture content: 10 questions, 4 variants each, correct variant in separate field.
func (c *Client) GenerateTest(ctx context.Context, lectureContent string) (*Test, error) {
	prompt := fmt.Sprintf(`Based on the following lecture content, generate exactly 10 multiple-choice test questions.

For each question:
1. Write a clear question that checks understanding of the lecture.
2. Provide exactly 4 answer variants (options A, B, C, D).
3. Include the correct answer text in a separate "correct_variant" field (must match one of the 4 variants exactly).

Lecture content:
---
%s
---

Respond with a valid JSON object only, no other text. Use this exact structure:
{"questions":[{"question":"...","variants":["A option","B option","C option","D option"],"correct_variant":"exact text of correct option"}]}`, lectureContent)

	reqBody := cohereRequest{
		Message: prompt,
		Model:   c.model,
		Stream:  false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var genResp cohereResponse
	if err := json.Unmarshal(body, &genResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if genResp.Text == "" {
		return nil, fmt.Errorf("no content generated")
	}

	var test Test
	if err := json.Unmarshal([]byte(genResp.Text), &test); err != nil {
		return nil, fmt.Errorf("failed to parse test JSON: %w", err)
	}

	if len(test.Questions) != 10 {
		return nil, fmt.Errorf("expected 10 questions, got %d", len(test.Questions))
	}

	for i, q := range test.Questions {
		if len(q.Variants) != 4 {
			return nil, fmt.Errorf("question %d: expected 4 variants, got %d", i+1, len(q.Variants))
		}
	}

	return &test, nil
}

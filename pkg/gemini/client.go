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

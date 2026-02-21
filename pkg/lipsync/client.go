package lipsync

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	client  *http.Client
}

type SyncResponse struct {
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type JobStatusResponse struct {
	JobID     string `json:"job_id"`
	Status    string `json:"status"`
	OutputURL string `json:"output_url"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 120 * time.Second, // allow time for multipart upload
		},
	}
}

// CreateLipsyncJobWithFiles sends avatar and audio as multipart/form-data body.
func (c *Client) CreateLipsyncJobWithFiles(ctx context.Context, avatarData, audioData []byte, avatarFilename, audioFilename string) (*SyncResponse, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	avatarPart, err := w.CreateFormFile("avatar", avatarFilename)
	if err != nil {
		return nil, fmt.Errorf("create avatar form file: %w", err)
	}
	if _, err := avatarPart.Write(avatarData); err != nil {
		return nil, fmt.Errorf("write avatar: %w", err)
	}

	audioPart, err := w.CreateFormFile("audio", audioFilename)
	if err != nil {
		return nil, fmt.Errorf("create audio form file: %w", err)
	}
	if _, err := audioPart.Write(audioData); err != nil {
		return nil, fmt.Errorf("write audio: %w", err)
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	url := fmt.Sprintf("%s/sync", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var syncResp SyncResponse
	if err := json.Unmarshal(body, &syncResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &syncResp, nil
}

func (c *Client) GetJobStatus(ctx context.Context, jobID string) (*JobStatusResponse, error) {
	url := fmt.Sprintf("%s/sync/%s", c.baseURL, jobID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

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

	var statusResp JobStatusResponse
	if err := json.Unmarshal(body, &statusResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &statusResp, nil
}

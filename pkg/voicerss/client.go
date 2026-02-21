package voicerss

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

type TTSRequest struct {
	Text     string
	Language string
	Voice    string
	Rate     int
	Codec    string
	Format   string
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://api.voicerss.org/",
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (c *Client) TextToSpeech(ctx context.Context, req TTSRequest) ([]byte, error) {
	if req.Language == "" {
		req.Language = "en-us"
	}
	if req.Codec == "" {
		req.Codec = "MP3"
	}
	if req.Format == "" {
		req.Format = "16khz_16bit_stereo"
	}

	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("src", req.Text)
	params.Set("hl", req.Language)
	params.Set("c", req.Codec)
	params.Set("f", req.Format)

	if req.Voice != "" {
		params.Set("v", req.Voice)
	}
	if req.Rate != 0 {
		params.Set("r", fmt.Sprintf("%d", req.Rate))
	}

	fullURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	httpReq, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(httpReq)
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

	contentType := resp.Header.Get("Content-Type")
	if contentType == "text/plain" {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	return body, nil
}

type Language struct {
	Code string
	Name string
}

func GetSupportedLanguages() []Language {
	return []Language{
		{Code: "en-us", Name: "English (United States)"},
		{Code: "en-gb", Name: "English (United Kingdom)"},
		{Code: "es-es", Name: "Spanish (Spain)"},
		{Code: "fr-fr", Name: "French (France)"},
		{Code: "de-de", Name: "German (Germany)"},
		{Code: "it-it", Name: "Italian (Italy)"},
		{Code: "pt-br", Name: "Portuguese (Brazil)"},
		{Code: "ru-ru", Name: "Russian (Russia)"},
		{Code: "zh-cn", Name: "Chinese (China)"},
		{Code: "ja-jp", Name: "Japanese (Japan)"},
		{Code: "ko-kr", Name: "Korean (South Korea)"},
		{Code: "ar-sa", Name: "Arabic (Saudi Arabia)"},
	}
}

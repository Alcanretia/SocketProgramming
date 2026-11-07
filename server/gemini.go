package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	defaultGeminiModel = "gemini-3.5-flash-lite"
	defaultGeminiBase  = "https://generativelanguage.googleapis.com/v1beta"
	maxGeminiResponse  = 1 << 20 // never read more than 1 MiB back
)

// NewGeminiSummarizer builds a summarizer that calls the Google Gemini API over
// HTTPS from the server (plain net/http, no extra dependency). The key lives
// only in the server's environment (GEMINI_API_KEY, or GOOGLE_API_KEY); it is
// sent in a header, never in the URL, and never reaches the browser. It returns
// nil when no key is configured.
//
// Env: GEMINI_MODEL overrides the model, GEMINI_BASE_URL overrides the endpoint
// (used by the tests to point at a fake server).
func NewGeminiSummarizer() Summarizer {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		key = os.Getenv("GOOGLE_API_KEY")
	}
	if key == "" {
		return nil
	}
	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = defaultGeminiModel
	}
	base := os.Getenv("GEMINI_BASE_URL")
	if base == "" {
		base = defaultGeminiBase
	}
	return newGeminiSummarizer(http.DefaultClient, base, model, key)
}

type geminiPart struct {
	Text    string `json:"text"`
	Thought bool   `json:"thought,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiRequest struct {
	SystemInstruction geminiContent   `json:"systemInstruction"`
	Contents          []geminiContent `json:"contents"`
	GenerationConfig  struct {
		MaxOutputTokens int     `json:"maxOutputTokens"`
		Temperature     float64 `json:"temperature"`
	} `json:"generationConfig"`
}

type geminiResponse struct {
	Candidates []struct {
		Content      geminiContent `json:"content"`
		FinishReason string        `json:"finishReason"`
	} `json:"candidates"`
	PromptFeedback struct {
		BlockReason string `json:"blockReason"`
	} `json:"promptFeedback"`
	Error *struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"error"`
}

func newGeminiSummarizer(client *http.Client, base, model, key string) Summarizer {
	endpoint := strings.TrimRight(base, "/") + "/models/" + url.PathEscape(model) + ":generateContent"

	return func(ctx context.Context, roomLabel string, lines []string) (string, error) {
		var body geminiRequest
		body.SystemInstruction = geminiContent{Parts: []geminiPart{{Text: summarySystemPrompt}}}
		body.Contents = []geminiContent{{Role: "user", Parts: []geminiPart{{Text: buildPrompt(roomLabel, lines)}}}}
		body.GenerationConfig.MaxOutputTokens = 2048 // headroom in case the model "thinks" first
		body.GenerationConfig.Temperature = 0.2      // a summary should be repeatable, not creative
		raw, err := json.Marshal(body)
		if err != nil {
			return "", err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(raw))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-goog-api-key", key)

		resp, err := client.Do(req)
		if err != nil {
			// The ctx error (timeout) stays inside the chain so the hub can tell
			// a timeout apart from other failures. *url.Error never contains the
			// header, so the key cannot leak through here.
			return "", err
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(io.LimitReader(resp.Body, maxGeminiResponse))
		if err != nil {
			return "", err
		}

		var out geminiResponse
		jsonErr := json.Unmarshal(data, &out)
		if resp.StatusCode != http.StatusOK {
			if jsonErr == nil && out.Error != nil {
				return "", fmt.Errorf("gemini: http %d %s: %s", resp.StatusCode, out.Error.Status, out.Error.Message)
			}
			return "", fmt.Errorf("gemini: http %d", resp.StatusCode)
		}
		if jsonErr != nil {
			return "", fmt.Errorf("gemini: unreadable response: %w", jsonErr)
		}
		if r := out.PromptFeedback.BlockReason; r != "" {
			return "", fmt.Errorf("gemini: the prompt was blocked (%s)", r)
		}
		if len(out.Candidates) == 0 {
			return "", errors.New("gemini: no candidates in the response")
		}
		var sb strings.Builder
		for _, p := range out.Candidates[0].Content.Parts {
			if !p.Thought { // skip the model's internal reasoning parts
				sb.WriteString(p.Text)
			}
		}
		text := strings.TrimSpace(sb.String())
		if text == "" {
			return "", fmt.Errorf("gemini: empty summary (finishReason=%s)", out.Candidates[0].FinishReason)
		}
		return text, nil
	}
}

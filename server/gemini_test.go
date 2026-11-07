package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// fakeGemini answers like the real endpoint and records what it was sent.
type fakeGemini struct {
	status int
	body   string
	delay  time.Duration

	gotPath, gotKey, gotCT string
	gotBody                geminiRequest
}

func (f *fakeGemini) start(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.gotPath, f.gotKey, f.gotCT = r.URL.Path, r.Header.Get("x-goog-api-key"), r.Header.Get("Content-Type")
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &f.gotBody)
		if f.delay > 0 {
			select {
			case <-time.After(f.delay):
			case <-r.Context().Done():
				return
			}
		}
		w.WriteHeader(f.status)
		_, _ = w.Write([]byte(f.body))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestGeminiSummarizerSuccess(t *testing.T) {
	f := &fakeGemini{status: 200, body: `{"candidates":[{"content":{"parts":[
		{"text":"hidden reasoning","thought":true},
		{"text":"- Bob said hi\n"},{"text":"- Bob asked a question"}]},"finishReason":"STOP"}]}`}
	srv := f.start(t)
	sum := newGeminiSummarizer(srv.Client(), srv.URL+"/v1beta", "test-model", "SECRET-KEY")

	got, err := sum(context.Background(), "DM with Bob", []string{"Bob: hi", "Bob: [sticker: smile]"})
	if err != nil {
		t.Fatal(err)
	}
	if want := "- Bob said hi\n- Bob asked a question"; got != want {
		t.Fatalf("got %q, want %q (thought parts must be skipped)", got, want)
	}
	if f.gotPath != "/v1beta/models/test-model:generateContent" {
		t.Fatalf("wrong path %q", f.gotPath)
	}
	if f.gotKey != "SECRET-KEY" || f.gotCT != "application/json" {
		t.Fatalf("key header %q, content-type %q", f.gotKey, f.gotCT)
	}
	if !strings.Contains(f.gotBody.SystemInstruction.Parts[0].Text, "untrusted user content") {
		t.Fatal("system prompt not sent as systemInstruction")
	}
	user := f.gotBody.Contents[0].Parts[0].Text
	if !strings.Contains(user, "<messages>") || !strings.Contains(user, "Bob: [sticker: smile]") || !strings.Contains(user, "DM with Bob") {
		t.Fatalf("user prompt missing wrapped messages: %q", user)
	}
}

func TestGeminiSummarizerFailures(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		want   string
	}{
		{"quota", 429, `{"error":{"code":429,"status":"RESOURCE_EXHAUSTED","message":"quota"}}`, "RESOURCE_EXHAUSTED"},
		{"bad key", 400, `{"error":{"code":400,"status":"INVALID_ARGUMENT","message":"API key not valid"}}`, "INVALID_ARGUMENT"},
		{"non-json error", 502, `<html>bad gateway</html>`, "http 502"},
		{"blocked", 200, `{"promptFeedback":{"blockReason":"SAFETY"}}`, "blocked"},
		{"no candidates", 200, `{}`, "no candidates"},
		{"empty text", 200, `{"candidates":[{"content":{"parts":[]},"finishReason":"MAX_TOKENS"}]}`, "MAX_TOKENS"},
		{"garbage", 200, `not json`, "unreadable"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := &fakeGemini{status: tc.status, body: tc.body}
			srv := f.start(t)
			sum := newGeminiSummarizer(srv.Client(), srv.URL, "m", "SECRET-KEY")
			_, err := sum(context.Background(), "room", []string{"a: b"})
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("want error containing %q, got %v", tc.want, err)
			}
			if strings.Contains(err.Error(), "SECRET-KEY") {
				t.Fatal("the API key must never appear in an error")
			}
		})
	}
}

func TestGeminiSummarizerTimeoutKeepsDeadlineError(t *testing.T) {
	f := &fakeGemini{status: 200, body: `{}`, delay: 2 * time.Second}
	srv := f.start(t)
	sum := newGeminiSummarizer(srv.Client(), srv.URL, "m", "SECRET-KEY")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := sum(ctx, "room", []string{"a: b"})
	// The hub maps this to the friendly "timeout" code.
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("want DeadlineExceeded in the chain, got %v", err)
	}
	if strings.Contains(err.Error(), "SECRET-KEY") {
		t.Fatal("the API key must never appear in an error")
	}
}

func TestGeminiKeyFromEnvironment(t *testing.T) {
	clear := func() {
		t.Setenv("GEMINI_API_KEY", "")
		t.Setenv("GOOGLE_API_KEY", "")
	}

	clear()
	if NewGeminiSummarizer() != nil {
		t.Fatal("no key: the summarizer must be nil so the server runs without AI")
	}

	clear()
	t.Setenv("GEMINI_API_KEY", "x")
	if NewGeminiSummarizer() == nil {
		t.Fatal("GEMINI_API_KEY should enable the summarizer")
	}

	clear()
	t.Setenv("GOOGLE_API_KEY", "x")
	if NewGeminiSummarizer() == nil {
		t.Fatal("GOOGLE_API_KEY should also enable the summarizer")
	}
}

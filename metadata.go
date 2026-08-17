package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// VideoInfo holds parsed metadata from yt-dlp --dump-json.
type VideoInfo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Uploader  string `json:"uploader"`
	Duration  int    `json:"duration"`
	Thumbnail string `json:"thumbnail"`
	WebpageURL string `json:"webpage_url"`
}

// FetchMetadata runs yt-dlp --dump-json for a single URL and returns parsed info.
func FetchMetadata(url string, verbose bool) (*VideoInfo, error) {
	done := make(chan struct{})
	go showSpinner("Fetching metadata", done)

	cmd := exec.Command("yt-dlp", "--dump-json", "--no-warnings", "--no-playlist", url)
	Verbose(verbose, "Running: %s", strings.Join(cmd.Args, " "))

	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	out, err := cmd.Output()
	close(done)
	ClearLine()

	if err != nil {
		errMsg := strings.TrimSpace(stderrBuf.String())
		if errMsg != "" {
			return nil, fmt.Errorf("failed to fetch metadata: %s", errMsg)
		}
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}

	var info VideoInfo
	if err := json.Unmarshal(out, &info); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}
	return &info, nil
}

func showSpinner(label string, done <-chan struct{}) {
	tick := 0
	for {
		select {
		case <-done:
			return
		default:
			ClearLine()
			fmt.Printf("  %s %s", SpinnerFrame(tick), label)
			tick++
			time.Sleep(80 * time.Millisecond)
		}
	}
}

// IsValidYouTubeURL checks if the provided URL is a valid YouTube URL.
func IsValidYouTubeURL(url string) bool {
	patterns := []string{
		`^https?://(www\.)?youtube\.com/watch\?v=[\w-]+`,
		`^https?://youtu\.be/[\w-]+`,
		`^https?://(www\.)?youtube\.com/embed/[\w-]+`,
	}
	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, url)
		if matched {
			return true
		}
	}
	return false
}

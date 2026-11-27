package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// VideoInfo holds parsed metadata from yt-dlp --dump-json.
type VideoInfo struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Uploader   string `json:"uploader"`
	Duration   int    `json:"duration"`
	Thumbnail  string `json:"thumbnail"`
	WebpageURL string `json:"webpage_url"`
}

// FetchMetadata runs yt-dlp --dump-json for a single URL and returns parsed info.
func FetchMetadata(url string, verbose bool) (*VideoInfo, error) {
	spin := StartSpinner("Fetching metadata")

	cmd := exec.Command("yt-dlp", "--dump-json", "--no-warnings", "--no-playlist", url)
	Verbose(verbose, "Running: %s", strings.Join(cmd.Args, " "))

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	spin.Stop()

	if err != nil {
		if msg := ytdlpError(stderr.String()); msg != "" {
			return nil, fmt.Errorf("couldn't get video info: %s", msg)
		}
		return nil, fmt.Errorf("couldn't get video info: %w", err)
	}

	var info VideoInfo
	if err := json.Unmarshal(out, &info); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}
	return &info, nil
}

// Spinner animates a label on the current line until Stop is called.
type Spinner struct {
	mu    sync.Mutex
	label string
	done  chan struct{}
	exit  chan struct{}
}

// StartSpinner begins animating label on the current line.
func StartSpinner(label string) *Spinner {
	s := &Spinner{label: label, done: make(chan struct{}), exit: make(chan struct{})}
	go s.run()
	return s
}

func (s *Spinner) run() {
	defer close(s.exit)
	t := time.NewTicker(80 * time.Millisecond)
	defer t.Stop()
	for tick := 0; ; tick++ {
		s.mu.Lock()
		label := s.label
		s.mu.Unlock()
		ClearLine()
		fmt.Printf("  %s %s", SpinnerFrame(tick), label)
		select {
		case <-s.done:
			return
		case <-t.C:
		}
	}
}

// SetLabel changes the text shown next to the spinner.
func (s *Spinner) SetLabel(label string) {
	s.mu.Lock()
	s.label = label
	s.mu.Unlock()
}

// Stop halts the animation, waits for the last frame, and clears the line.
// It is safe to call more than once.
func (s *Spinner) Stop() {
	select {
	case <-s.done:
	default:
		close(s.done)
	}
	<-s.exit
	ClearLine()
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Options holds all CLI-derived settings.
type Options struct {
	Quality   int
	OutputDir string
	KeepVideo bool
	NoMeta    bool
	Verbose   bool
}

// CheckDependencies verifies that yt-dlp and ffmpeg are installed.
func CheckDependencies() error {
	missing := []string{}

	if _, err := exec.LookPath("yt-dlp"); err != nil {
		missing = append(missing, "yt-dlp")
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		missing = append(missing, "ffmpeg")
	}

	if len(missing) > 0 {
		msg := "Missing required dependencies: " + strings.Join(missing, ", ") + "\n\n"
		msg += "  Install with:\n"
		for _, dep := range missing {
			switch dep {
			case "yt-dlp":
				msg += "    brew install yt-dlp       (macOS)\n"
				msg += "    pip install yt-dlp         (pip)\n"
				msg += "    sudo apt install yt-dlp    (Debian/Ubuntu)\n"
			case "ffmpeg":
				msg += "    brew install ffmpeg        (macOS)\n"
				msg += "    sudo apt install ffmpeg    (Debian/Ubuntu)\n"
			}
		}
		return fmt.Errorf("%s", msg)
	}
	return nil
}

// DownloadAndConvert downloads a single URL and converts it to MP3.
func DownloadAndConvert(url string, opts Options) error {
	info, err := FetchMetadata(url, opts.Verbose)
	if err != nil {
		return err
	}
	PrintTrackInfo(info)

	Info("Downloading and converting to MP3 (%d kbps)…", opts.Quality)

	// yt-dlp sanitizes titles into file names (e.g. "/" and ":" are replaced),
	// so ask it to write the final path instead of guessing it.
	pathFile, err := os.CreateTemp("", "konversion-path-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	pathFile.Close()
	defer os.Remove(pathFile.Name())

	args := []string{
		url,
		"--no-playlist", // a watch?v=…&list=… link must not pull the whole playlist
		"--extract-audio",
		"--audio-format", "mp3",
		"--audio-quality", strconv.Itoa(opts.Quality) + "K",
		"--newline", // one progress line per update
		"-o", filepath.Join(opts.OutputDir, "%(title)s.%(ext)s"),
		"--print-to-file", "after_move:filepath", pathFile.Name(),
	}

	if !opts.NoMeta {
		args = append(args, "--embed-thumbnail", "--embed-metadata")
		// yt-dlp needs atomicparsley or mutagen for thumbnails on mp3;
		// --convert-thumbnails jpg ensures compatibility
		args = append(args, "--convert-thumbnails", "jpg")
	}

	if opts.KeepVideo {
		args = append(args, "-k")
	}

	Verbose(opts.Verbose, "Running: yt-dlp %s", strings.Join(args, " "))

	cmd := exec.Command("yt-dlp", args...)

	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to capture output: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start yt-dlp: %w", err)
	}

	var spin *Spinner
	stopSpin := func() {
		if spin != nil {
			spin.Stop()
			spin = nil
		}
	}
	phase := func(label string) {
		if spin == nil {
			ClearLine()
			spin = StartSpinner(label)
		} else {
			spin.SetLabel(label)
		}
	}

	scanner := bufio.NewScanner(stdout)
	progressRe := regexp.MustCompile(`^\[download\]\s+([\d.]+)%\s+of`)
	best := 0.0
	for scanner.Scan() {
		line := scanner.Text()
		if opts.Verbose {
			stopSpin()
			Verbose(true, "%s", line)
		}

		switch {
		case progressRe.MatchString(line):
			if spin != nil {
				continue // post-processing already started; ignore stray downloads (thumbnails)
			}
			pct, _ := strconv.ParseFloat(progressRe.FindStringSubmatch(line)[1], 64)
			// Never move the bar backwards when yt-dlp restarts a fragment or stream.
			if pct > best {
				best = pct
			}
			if !opts.Verbose {
				ProgressBar(best, 30, "")
			}
		case strings.HasPrefix(line, "[ExtractAudio]"):
			if !opts.Verbose {
				phase("Converting to MP3")
			}
		case strings.HasPrefix(line, "[Metadata]"), strings.HasPrefix(line, "[EmbedThumbnail]"),
			strings.HasPrefix(line, "[ThumbnailsConvertor]"):
			if !opts.Verbose && spin != nil {
				phase("Adding title, artist and cover art")
			}
		}
	}
	stopSpin()

	if err := cmd.Wait(); err != nil {
		ClearLine()
		if msg := ytdlpError(stderrBuf.String()); msg != "" {
			return fmt.Errorf("yt-dlp failed: %s", msg)
		}
		return fmt.Errorf("yt-dlp failed: %w", err)
	}

	ClearLine()

	// Show full file path and size
	savedPath := lastLine(pathFile.Name())
	if savedPath == "" {
		fmt.Printf("Saved to: %s\n", opts.OutputDir)
		return nil
	}
	if fi, err := os.Stat(savedPath); err == nil {
		sizeMB := float64(fi.Size()) / 1024 / 1024
		fmt.Printf("Saved: %s (%.1f MB)\n", savedPath, sizeMB)
	} else {
		fmt.Printf("Saved: %s\n", savedPath)
	}
	return nil
}

// lastLine returns the last non-empty line of a file, or "" if unreadable.
func lastLine(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	return strings.TrimSpace(lines[len(lines)-1])
}

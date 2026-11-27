package main

import (
	"net/url"
	"regexp"
	"strings"
)

var videoIDRe = regexp.MustCompile(`^[\w-]{11}$`)

// NormalizeYouTubeURL accepts the common forms of a YouTube video link and
// returns a canonical https://www.youtube.com/watch?v=ID URL with no playlist
// or tracking parameters. ok is false if the input isn't a single-video link.
//
// Accepted: youtube.com/watch?v=, youtu.be/, /shorts/, /embed/, /live/,
// /v/, on the www., m. and music. hosts, with or without the scheme.
func NormalizeYouTubeURL(raw string) (canonical string, ok bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", false
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return "", false
	}

	host := strings.ToLower(u.Hostname())
	path := strings.Trim(u.Path, "/")
	var id string

	switch host {
	case "youtu.be", "www.youtu.be":
		id = firstSegment(path)
	case "youtube.com", "www.youtube.com", "m.youtube.com", "music.youtube.com",
		"youtube-nocookie.com", "www.youtube-nocookie.com":
		if path == "watch" {
			id = u.Query().Get("v")
			break
		}
		for _, prefix := range []string{"shorts/", "embed/", "live/", "v/"} {
			if strings.HasPrefix(path, prefix) {
				id = firstSegment(strings.TrimPrefix(path, prefix))
				break
			}
		}
	default:
		return "", false
	}

	if !videoIDRe.MatchString(id) {
		return "", false
	}
	return "https://www.youtube.com/watch?v=" + id, true
}

func firstSegment(path string) string {
	if i := strings.IndexByte(path, '/'); i >= 0 {
		return path[:i]
	}
	return path
}

// ytdlpError turns yt-dlp's stderr into a short, readable reason. It prefers
// the ERROR: lines and drops yt-dlp's "[youtube] ID:" prefix noise.
func ytdlpError(stderr string) string {
	var errs []string
	for _, line := range strings.Split(stderr, "\n") {
		line = strings.TrimSpace(line)
		if msg, found := strings.CutPrefix(line, "ERROR: "); found {
			// "[youtube] dQw4w9WgXcQ: Video unavailable" -> "Video unavailable"
			if strings.HasPrefix(msg, "[") {
				if i := strings.Index(msg, ": "); i >= 0 {
					msg = msg[i+2:]
				}
			}
			errs = append(errs, msg)
		}
	}
	if len(errs) > 0 {
		return strings.Join(errs, "; ")
	}
	return strings.TrimSpace(stderr)
}

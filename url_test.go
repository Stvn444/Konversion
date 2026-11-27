package main

import "testing"

func TestNormalizeYouTubeURL(t *testing.T) {
	const want = "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	good := []string{
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		"http://youtube.com/watch?v=dQw4w9WgXcQ",
		"www.youtube.com/watch?v=dQw4w9WgXcQ",
		"https://m.youtube.com/watch?v=dQw4w9WgXcQ",
		"https://music.youtube.com/watch?v=dQw4w9WgXcQ&feature=share",
		"https://www.youtube.com/watch?feature=share&v=dQw4w9WgXcQ",
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ&list=PL1234567890&index=3",
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ&t=42s",
		"https://youtu.be/dQw4w9WgXcQ",
		"https://youtu.be/dQw4w9WgXcQ?si=abcdef&t=10",
		"https://www.youtube.com/shorts/dQw4w9WgXcQ",
		"https://www.youtube.com/embed/dQw4w9WgXcQ?start=5",
		"https://www.youtube.com/live/dQw4w9WgXcQ",
		"  https://youtu.be/dQw4w9WgXcQ  ",
	}
	for _, in := range good {
		got, ok := NormalizeYouTubeURL(in)
		if !ok || got != want {
			t.Errorf("NormalizeYouTubeURL(%q) = %q, %v; want %q, true", in, got, ok, want)
		}
	}

	bad := []string{
		"",
		"quit",
		"https://www.youtube.com/playlist?list=PL1234567890",
		"https://www.youtube.com/watch",
		"https://www.youtube.com/watch?v=short",
		"https://www.youtube.com/@somechannel",
		"https://evilyoutube.com/watch?v=dQw4w9WgXcQ",
		"https://youtube.com.evil.example/watch?v=dQw4w9WgXcQ",
		"ftp://youtube.com/watch?v=dQw4w9WgXcQ",
		"https://vimeo.com/123456",
	}
	for _, in := range bad {
		if got, ok := NormalizeYouTubeURL(in); ok {
			t.Errorf("NormalizeYouTubeURL(%q) = %q, true; want rejected", in, got)
		}
	}
}

func TestYtdlpError(t *testing.T) {
	stderr := "WARNING: something minor\nERROR: [youtube] dQw4w9WgXcQ: Video unavailable. This video is private\n"
	if got, want := ytdlpError(stderr), "Video unavailable. This video is private"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if got, want := ytdlpError("  plain failure \n"), "plain failure"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

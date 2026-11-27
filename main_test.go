package main

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	home, _ := os.UserHomeDir()
	dl := filepath.Join(home, "Downloads")

	cases := []struct {
		args []string
		opts Options
		urls []string
	}{
		{nil, Options{Quality: 320, OutputDir: dl}, nil},
		{[]string{"-q", "192", "URL1"}, Options{Quality: 192, OutputDir: dl}, []string{"URL1"}},
		{[]string{"URL1", "--quality=128", "URL2", "-k"},
			Options{Quality: 128, OutputDir: dl, KeepVideo: true}, []string{"URL1", "URL2"}},
		{[]string{"-o", "~/Music/samples", "--no-meta", "-v", "URL1"},
			Options{Quality: 320, OutputDir: filepath.Join(home, "Music/samples"), NoMeta: true, Verbose: true},
			[]string{"URL1"}},
		{[]string{"--output", "/tmp/x", "--keep-video", "--verbose"},
			Options{Quality: 320, OutputDir: "/tmp/x", KeepVideo: true, Verbose: true}, nil},
	}
	for _, c := range cases {
		opts, urls, _, err := parseArgs(c.args, io.Discard)
		if err != nil {
			t.Errorf("parseArgs(%q): %v", c.args, err)
			continue
		}
		if opts != c.opts || !reflect.DeepEqual(urls, c.urls) {
			t.Errorf("parseArgs(%q) = %+v %q; want %+v %q", c.args, opts, urls, c.opts, c.urls)
		}
	}

	for _, bad := range [][]string{{"-q", "300"}, {"-q", "abc"}, {"--nope"}} {
		if _, _, _, err := parseArgs(bad, io.Discard); err == nil {
			t.Errorf("parseArgs(%q) succeeded; want error", bad)
		}
	}

	if _, _, v, err := parseArgs([]string{"--version"}, io.Discard); err != nil || !v {
		t.Errorf("--version not detected: %v %v", v, err)
	}
}

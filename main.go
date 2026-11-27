package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

const usage = `Usage: konversion [options] [URL ...]

With no URLs, Konversion starts an interactive prompt: paste links one at a time.
With URLs, it converts each one and exits.

Options:
  -q, --quality KBPS   MP3 bitrate: 128, 192, 256 or 320 (default 320)
  -o, --output DIR     folder to save MP3s in (default ~/Downloads)
  -k, --keep-video     also keep the full video (downloads video, not just audio)
      --no-meta        don't embed title, artist or cover art
  -v, --verbose        show the yt-dlp commands and their output
      --version        print the version and exit
  -h, --help           show this help
`

var validQualities = map[int]bool{128: true, 192: true, 256: true, 320: true}

// parseArgs reads options and URLs. Options may come before or after URLs.
func parseArgs(args []string, stderr io.Writer) (opts Options, urls []string, showVersion bool, err error) {
	home, _ := os.UserHomeDir()
	opts = Options{Quality: 320, OutputDir: filepath.Join(home, "Downloads")}

	fs := flag.NewFlagSet("konversion", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprint(stderr, usage) }

	fs.IntVar(&opts.Quality, "q", opts.Quality, "")
	fs.IntVar(&opts.Quality, "quality", opts.Quality, "")
	fs.StringVar(&opts.OutputDir, "o", opts.OutputDir, "")
	fs.StringVar(&opts.OutputDir, "output", opts.OutputDir, "")
	fs.BoolVar(&opts.KeepVideo, "k", false, "")
	fs.BoolVar(&opts.KeepVideo, "keep-video", false, "")
	fs.BoolVar(&opts.NoMeta, "no-meta", false, "")
	fs.BoolVar(&opts.Verbose, "v", false, "")
	fs.BoolVar(&opts.Verbose, "verbose", false, "")
	fs.BoolVar(&showVersion, "version", false, "")

	// The flag package stops at the first non-flag argument, so keep parsing
	// after each URL to allow `konversion URL -q 192`.
	for {
		if err = fs.Parse(args); err != nil {
			return
		}
		if fs.NArg() == 0 {
			break
		}
		urls = append(urls, fs.Arg(0))
		args = fs.Args()[1:]
	}

	if !validQualities[opts.Quality] {
		err = fmt.Errorf("invalid quality %d: use 128, 192, 256 or 320", opts.Quality)
		return
	}
	opts.OutputDir = expandHome(opts.OutputDir, home)
	return
}

func expandHome(path, home string) string {
	if path == "~" {
		return home
	}
	if rest, ok := strings.CutPrefix(path, "~/"); ok {
		return filepath.Join(home, rest)
	}
	return path
}

func main() {
	opts, urls, showVersion, err := parseArgs(os.Args[1:], os.Stderr)
	if err == flag.ErrHelp {
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n%s", err, usage)
		os.Exit(2)
	}
	if showVersion {
		fmt.Printf("konversion %s\n", version)
		return
	}

	// Check dependencies
	if len(urls) == 0 {
		PrintBanner()
		fmt.Println("Checking dependencies...")
	}
	if err := CheckDependencies(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		fmt.Printf("Error: can't create output folder: %v\n", err)
		os.Exit(1)
	}

	if len(urls) > 0 {
		os.Exit(runBatch(urls, opts))
	}

	fmt.Println("Dependencies OK!")
	fmt.Println()
	runInteractive(opts)
}

// runBatch converts each URL from the command line and returns the exit code:
// 0 if all succeeded, 1 if any failed.
func runBatch(urls []string, opts Options) int {
	failed := 0
	for _, input := range urls {
		if !convertOne(input, opts) {
			failed++
		}
	}
	if len(urls) > 1 {
		fmt.Printf("\nDone: %d of %d converted.\n", len(urls)-failed, len(urls))
	}
	if failed > 0 {
		return 1
	}
	return 0
}

func runInteractive(opts Options) {
	for {
		fmt.Print("Enter YouTube URL (or 'quit' to exit): ")
		input, err := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Ctrl-D / closed stdin: stop instead of looping forever.
		if err != nil && input == "" {
			fmt.Println()
			fmt.Println("Goodbye!")
			return
		}

		if input == "quit" || input == "q" {
			fmt.Println("Goodbye!")
			return
		}

		if input == "" {
			continue
		}

		convertOne(input, opts)
		fmt.Println()
	}
}

// convertOne validates and converts a single link, printing any error.
func convertOne(input string, opts Options) bool {
	url, ok := NormalizeYouTubeURL(input)
	if !ok {
		fmt.Printf("Error: Invalid YouTube URL format: %s\n", input)
		fmt.Println("Supported formats:")
		fmt.Println("  - https://www.youtube.com/watch?v=VIDEO_ID")
		fmt.Println("  - https://youtu.be/VIDEO_ID")
		fmt.Println("  - https://www.youtube.com/shorts/VIDEO_ID")
		fmt.Println("  (playlist and channel links aren't supported)")
		fmt.Println()
		return false
	}

	if err := DownloadAndConvert(url, opts); err != nil {
		fmt.Printf("Error: %v\n", err)
		return false
	}
	return true
}

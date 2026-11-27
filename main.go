package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func main() {
	PrintBanner()

	// Check dependencies
	fmt.Println("Checking dependencies...")
	if err := CheckDependencies(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Dependencies OK!")
	fmt.Println()

	home, _ := os.UserHomeDir()
	downloadDir := filepath.Join(home, "Downloads")

	// Create Downloads folder if it doesn't exist
	os.MkdirAll(downloadDir, 0755)

	opts := Options{
		Quality:   320,
		OutputDir: downloadDir,
	}

	for {
		fmt.Print("Enter YouTube URL (or 'quit' to exit): ")
		input, err := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Ctrl-D / closed stdin: stop instead of looping forever.
		if err != nil && input == "" {
			fmt.Println()
			fmt.Println("Goodbye!")
			break
		}

		if input == "quit" || input == "q" {
			fmt.Println("Goodbye!")
			break
		}

		if input == "" {
			continue
		}

		// Validate URL
		url, ok := NormalizeYouTubeURL(input)
		if !ok {
			fmt.Println("Error: Invalid YouTube URL format")
			fmt.Println("Supported formats:")
			fmt.Println("  - https://www.youtube.com/watch?v=VIDEO_ID")
			fmt.Println("  - https://youtu.be/VIDEO_ID")
			fmt.Println("  - https://www.youtube.com/shorts/VIDEO_ID")
			fmt.Println("  (playlist and channel links aren't supported)")
			fmt.Println()
			continue
		}

		// Download and convert
		if err := DownloadAndConvert(url, opts); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Println()
	}
}

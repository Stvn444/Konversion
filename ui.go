package main

import (
	"fmt"
	"strings"
)

// ANSI color codes
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Magenta   = "\033[35m"
	Cyan      = "\033[36m"
	White     = "\033[37m"
	BoldCyan  = "\033[1;36m"
	BoldGreen = "\033[1;32m"
	BoldRed   = "\033[1;31m"
)

const banner = `
 ██╗  ██╗ ██████╗ ███╗   ██╗██╗   ██╗███████╗██████╗ ███████╗██╗ ██████╗ ███╗   ██╗
 ██║ ██╔╝██╔═══██╗████╗  ██║██║   ██║██╔════╝██╔══██╗██╔════╝██║██╔═══██╗████╗  ██║
 █████╔╝ ██║   ██║██╔██╗ ██║██║   ██║█████╗  ██████╔╝███████╗██║██║   ██║██╔██╗ ██║
 ██╔═██╗ ██║   ██║██║╚██╗██║╚██╗ ██╔╝██╔══╝  ██╔══██╗╚════██║██║██║   ██║██║╚██╗██║
 ██║  ██╗╚██████╔╝██║ ╚████║ ╚████╔╝ ███████╗██║  ██║███████║██║╚██████╔╝██║ ╚████║
 ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝╚══════╝╚═╝ ╚═════╝ ╚═╝  ╚═╝
`

const version = "1.1.0"

func PrintBanner() {
	fmt.Print(Cyan + banner + Reset)
	fmt.Printf("  %sYouTube to MP3%s  v%s\n", Dim, Reset, version)
	fmt.Printf("  %sConvert any YouTube video to high quality MP3%s\n\n", Dim, Reset)
}

func colorize(color, text string) string {
	return color + text + Reset
}

func Info(format string, args ...any) {
	fmt.Printf(colorize(Cyan, "  ℹ ")+format+"\n", args...)
}

func Success(format string, args ...any) {
	fmt.Printf(colorize(Green, "  ✓ ")+format+"\n", args...)
}

func Warn(format string, args ...any) {
	fmt.Printf(colorize(Yellow, "  ⚠ ")+format+"\n", args...)
}

func Error(format string, args ...any) {
	fmt.Printf(colorize(Red, "  ✗ ")+format+"\n", args...)
}

func Verbose(verbose bool, format string, args ...any) {
	if verbose {
		fmt.Printf(colorize(Dim, "  … ")+format+"\n", args...)
	}
}

// ProgressBar renders a terminal progress bar.
// percent should be 0.0–100.0.
func ProgressBar(percent float64, width int, extra string) {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := int(percent / 100 * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled

	bar := Green + strings.Repeat("█", filled) + Dim + strings.Repeat("░", empty) + Reset
	fmt.Printf("\r  %s %s%.1f%%%s %s", bar, Bold, percent, Reset, extra)
}

func ClearLine() {
	fmt.Print("\r\033[K")
}

// Spinner frames for metadata fetching
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func SpinnerFrame(tick int) string {
	return Cyan + spinnerFrames[tick%len(spinnerFrames)] + Reset
}

func FormatDuration(seconds int) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%d:%02d", m, s)
}

func PrintTrackInfo(info *VideoInfo) {
	fmt.Println()
	fmt.Printf("  %sTitle:%s    %s\n", Bold, Reset, info.Title)
	fmt.Printf("  %sChannel:%s  %s\n", Bold, Reset, info.Uploader)
	fmt.Printf("  %sDuration:%s %s\n", Bold, Reset, FormatDuration(info.Duration))
	fmt.Println()
}


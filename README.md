# Konversion

```
 ██╗  ██╗ ██████╗ ███╗   ██╗██╗   ██╗███████╗██████╗ ███████╗██╗ ██████╗ ███╗   ██╗
 ██║ ██╔╝██╔═══██╗████╗  ██║██║   ██║██╔════╝██╔══██╗██╔════╝██║██╔═══██╗████╗  ██║
 █████╔╝ ██║   ██║██╔██╗ ██║██║   ██║█████╗  ██████╔╝███████╗██║██║   ██║██╔██╗ ██║
 ██╔═██╗ ██║   ██║██║╚██╗██║╚██╗ ██╔╝██╔══╝  ██╔══██╗╚════██║██║██║   ██║██║╚██╗██║
 ██║  ██╗╚██████╔╝██║ ╚████║ ╚████╔╝ ███████╗██║  ██║███████║██║╚██████╔╝██║ ╚████║
 ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝╚══════╝╚═╝ ╚═════╝ ╚═╝  ╚═╝
```

**YouTube to MP3 converter for the terminal.** Paste a link, get an MP3. That's it.

Built for producers and artists who just want to rip audio from YouTube without dealing with sketchy websites or bloated apps.

---

## What It Does

- Paste any YouTube URL and it downloads + converts to high quality MP3 (320kbps)
- Shows you the track title, channel, and duration before downloading
- Progress bar so you can see the download happening
- Automatically embeds metadata (title, artist, thumbnail) into the MP3
- MP3s are saved to your Downloads folder
- Shows you the exact file path and size when it's done

---

## Quick Install (macOS / Linux)

Open **Terminal** and paste this single command:

```bash
curl -fsSL https://raw.githubusercontent.com/Stvn444/Konversion/main/install.sh | bash
```

This automatically downloads and installs Konversion for your system.

**You still need yt-dlp and ffmpeg.** The installer will tell you if they're missing. Install them:

**macOS:**

```bash
brew install yt-dlp ffmpeg
```

> Don't have Homebrew? Install it first: `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`

**Linux (Ubuntu/Debian):**

```bash
sudo apt install ffmpeg
```

```bash
pip install yt-dlp
```

Then just type `konversion` to run it.

---

## Install (Windows)

1. Download `konversion-windows-amd64.exe` from the [Releases page](https://github.com/Stvn444/Konversion/releases)
2. Rename it to `konversion.exe` and put it somewhere in your PATH
3. Install [ffmpeg](https://www.gyan.dev/ffmpeg/builds/) — download "ffmpeg-release-essentials.zip", extract it, and add the `bin` folder to your PATH
4. Install [yt-dlp](https://github.com/yt-dlp/yt-dlp/releases) — download `yt-dlp.exe` and put it in a folder that's in your PATH

Then open Command Prompt or PowerShell and type:

```bash
konversion
```

---

## How to Use

Just type `konversion` in Terminal (all lowercase):

```bash
konversion
```

You'll see the Konversion banner and a prompt. Paste your YouTube URL and press Enter:

```
Enter YouTube URL (or 'quit' to exit): https://www.youtube.com/watch?v=dQw4w9WgXcQ

  Title:    Rick Astley - Never Gonna Give You Up
  Channel:  Rick Astley
  Duration: 3:33

  Downloading and converting to MP3 (320 kbps)...
  ██████████████████████████████ 100.0%
Saved: /Users/you/Downloads/Rick Astley - Never Gonna Give You Up.mp3 (8.2 MB)

Enter YouTube URL (or 'quit' to exit):
```

It loops — keep pasting URLs back to back. Type `quit` when you're done.

### Where Do My MP3s Go?

Your **Downloads** folder. After every download it prints the full path so you know exactly where to find it.

---

## Troubleshooting

### "command not found: konversion"

Re-run the install script:

```bash
curl -fsSL https://raw.githubusercontent.com/Stvn444/Konversion/main/install.sh | bash
```

Make sure you're typing it **all lowercase**: `konversion` (not `Konversion`).

### "yt-dlp is not installed"

```bash
brew install yt-dlp         # macOS
pip install yt-dlp           # any system with Python
sudo apt install yt-dlp      # Ubuntu/Debian
```

### "ffmpeg is not installed"

```bash
brew install ffmpeg          # macOS
sudo apt install ffmpeg      # Ubuntu/Debian
```

### "Invalid YouTube URL format"

Make sure you're pasting the full URL. These formats work:

- `https://www.youtube.com/watch?v=VIDEO_ID`
- `https://youtu.be/VIDEO_ID`
- `https://www.youtube.com/shorts/VIDEO_ID`
- `https://www.youtube.com/embed/VIDEO_ID` and `/live/VIDEO_ID`
- `m.youtube.com` and `music.youtube.com` links work too

Playlist and channel links aren't supported. A video link that's part of a playlist (`watch?v=…&list=…`) downloads just that one video.

### How do I update Konversion?

Just run the install script again:

```bash
curl -fsSL https://raw.githubusercontent.com/Stvn444/Konversion/main/install.sh | bash
```

---

## How It Works

Konversion is a Go wrapper around two industry-standard open source tools:

- **[yt-dlp](https://github.com/yt-dlp/yt-dlp)** — downloads video/audio from YouTube
- **[ffmpeg](https://ffmpeg.org/)** — converts audio to MP3 format

Konversion handles the terminal UI, progress bar, metadata display, and makes the whole process as simple as paste-and-go. No Python scripts, no browser extensions, no sketchy websites with popups.

## Disclaimer

Konversion is provided for **personal and educational use only**. This tool is intended for downloading content that you have the right to download, such as:

- Videos you own or have created
- Content licensed under Creative Commons or similar open licenses
- Content where the creator has explicitly allowed downloading
- Fair use purposes (commentary, education, research)

**You are solely responsible for how you use this tool.** Downloading copyrighted material without permission may violate YouTube's Terms of Service and copyright laws in your country. The author(s) of Konversion are not responsible for any misuse of this software or any violations of third-party rights, terms of service, or applicable laws.

By using Konversion, you agree that you will comply with all applicable laws and respect the rights of content creators.

## License

MIT — do whatever you want with it.

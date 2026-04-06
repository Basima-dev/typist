```markdown
<div align="center">
<pre>
╔══════════════════════════════════════════════════════════════════╗
║                                                                  ║
║     ████████╗██╗   ██╗██████╗ ██╗███████╗████████╗              ║
║     ╚══██╔══╝╚██╗ ██╔╝██╔══██╗██║██╔════╝╚══██╔══╝              ║
║        ██║    ╚████╔╝ ██████╔╝██║███████╗   ██║                 ║
║        ██║     ╚██╔╝  ██╔═══╝ ██║╚════██║   ██║                 ║
║        ██║      ██║   ██║     ██║███████║   ██║                 ║
║        ╚═╝      ╚═╝   ╚═╝     ╚═╝╚══════╝   ╚═╝                 ║
║                                                                  ║
║              A fast, offline typing test                         ║
║         No account. No paywall. No internet required.            ║
║                                                                  ║
╚══════════════════════════════════════════════════════════════════╝
</pre>
<p>
  <a href="https://github.com/chuma-beep/typist/stargazers"><img src="https://img.shields.io/github/stars/chuma-beep/typist?style=flat-square&color=yellow&logo=github" alt="stars"></a>
  <a href="https://github.com/chuma-beep/typist/network/members"><img src="https://img.shields.io/github/forks/chuma-beep/typist?style=flat-square&color=blue&logo=github" alt="forks"></a>
  <a href="https://github.com/chuma-beep/typist/issues"><img src="https://img.shields.io/github/issues/chuma-beep/typist?style=flat-square&color=red&logo=github" alt="issues"></a>
  <a href="LICENSE"><img src="https://img.shields.io/github/license/chuma-beep/typist?style=flat-square&color=green&logo=open-source-initiative" alt="license"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white" alt="go version"></a>
</p>
<pre>
┌─────────────────────────────────────────────────────────────────┐
│  $ typist          # Terminal UI                                │
│  $ typist --web    # Web UI (auto-opens browser)                │
│  $ typist --help   # Show all options                           │
└─────────────────────────────────────────────────────────────────┘
</pre>
</div>

---

## Quick Start

```bash
# Clone and build
git clone https://github.com/chuma-beep/typist
cd typist && go mod tidy && go build -o typist .

# Run
./typist          # Terminal UI
./typist --web    # Web UI
```

---

## Features

<pre>
┌────────────────────────────────┬──────────┬──────────┐
│ Feature                        │ Terminal │   Web    │
├────────────────────────────────┼──────────┼──────────┤
│ Word Mode — 30 common words    │    ✓     │    ✓     │
│ Time Mode — 15/30/60/120s      │    ✓     │    ✓     │
│ Quote Mode — Literary excerpts │    ✓     │    ✓     │
│ Code Mode — Go/JS/Python/Rust  │    ✓     │    ✓     │
│ Syntax Highlighting (Chroma)   │    ✓     │    ✓     │
│ Live WPM + Accuracy Stats      │    ✓     │    ✓     │
│ WPM Graph Over Time            │ Sparkline│ Chart.js │
│ Mistake Heatmap                │ Top-6    │ Keyboard │
│ Blind Mode (muscle memory)     │    ✓     │    ✗     │
│ Persistent Personal Bests      │    ✓     │    ✓     │
│ Session History (last 200)     │    ✓     │    ✗     │
│ Export to JSON / CSV           │    ✓     │    ✗     │
│ Single Binary, Zero Deps       │    ✓     │    ✓     │
└────────────────────────────────┴──────────┴──────────┘
</pre>

---

## Terminal UI

<pre>
┌─────────────────────────────────────────────────────────────────┐
│  CONTROLS                                                       │
├─────────────────────────────────────────────────────────────────┤
│  Menu                                                           │
│    ← →     Switch mode                                          │
│    ↑ ↓     Switch sub-row (time/lang)                           │
│    Enter   Start test                                           │
│    Esc/q   Quit                                                 │
│                                                                 │
│  Typing                                                         │
│    Ctrl+R  Restart with new text                                │
│    Ctrl+B  Toggle Blind Mode                                    │
│    Tab     Type tab (code mode)                                 │
│    Enter   Type newline (code mode)                             │
│    Esc     Quit                                                 │
│                                                                 │
│  Results                                                        │
│    Enter/R  Try again                                           │
│    M        Back to menu                                        │
│    H        View session history                                │
│    J        Export to JSON                                      │
│    C        Export to CSV                                       │
│    Esc      Quit                                                │
└─────────────────────────────────────────────────────────────────┘
</pre>

---

## Code Mode

Type real snippets with syntax highlighting powered by **Chroma**:

<pre>
┌────────────┬──────────┬────────────────────────────────────────┐
│ Language   │ Snippets │ Examples                               │
├────────────┼──────────┼────────────────────────────────────────┤
│ Go         │    8     │ Generics, channels, linked lists       │
│ JavaScript │    6     │ Debounce, memoize, EventEmitter        │
│ Python     │    5     │ Quicksort, LRU cache, decorators       │
│ Rust       │    5     │ Pattern matching, traits, generics     │
└────────────┴──────────┴────────────────────────────────────────┘
</pre>

---

## Blind Mode

<pre>
┌─────────────────────────────────────────────────────────────────┐
│  Ctrl+B  →  Every char becomes · (green=correct, red=wrong)     │
│                                                                 │
│  Forces typing from memory. Essential for muscle memory.        │
└─────────────────────────────────────────────────────────────────┘
</pre>

---

## Scores & Export

<pre>
┌─────────────────────────────────────────────────────────────────┐
│  Storage:  ~/.typist/scores.json                               │
│                                                                 │
│  Export (from results screen):                                  │
│    J  →  ~/typist-export-<timestamp>.json                       │
│    C  →  ~/typist-export-<timestamp>.csv                        │
└─────────────────────────────────────────────────────────────────┘
</pre>

---

## Tech Stack

<pre>
┌────────────────────┬─────────────────────────┬──────────────────┐
│ Component          │ Technology              │ Purpose          │
├────────────────────┼─────────────────────────┼──────────────────┤
│ TUI Framework      │ Bubble Tea              │ Elm architecture │
│ Terminal Styling   │ Lipgloss                │ Colors/layout    │
│ Syntax Highlight   │ Chroma v2               │ 300+ languages   │
│ Web Charts         │ Chart.js                │ WPM graphs       │
│ Core               │ Go Standard Library     │ HTTP, JSON, CSV  │
└────────────────────┴─────────────────────────┴──────────────────┘
</pre>

---

## Architecture

```
typist/
├── main.go          # Entry point, --web flag
├── model.go         # Bubble Tea model
├── highlight.go     # Chroma → lipgloss
├── words.go         # Text generation, wrapping
├── snippets.go      # Code library
├── scores.go        # Persistence, export
├── styles.go        # Catppuccin Mocha
├── web.go           # HTTP server
├── web/index.html   # Single-file web UI
└── quotes.json      # Embedded quotes
```

---

## Roadmap

- [ ] WPM sparkline → bar chart in TUI
- [ ] Dark/light theme toggle
- [ ] Focus mode (hide stats)
- [ ] Custom text input
- [ ] WebAssembly build

---

<pre>
┌─────────────────────────────────────────────────────────────────┐
│  MIT License  │  github.com/chuma-beep/typist                  │
└─────────────────────────────────────────────────────────────────┘
</pre>
```

package main

import "github.com/charmbracelet/lipgloss"

// Theme holds the full Catppuccin palette plus semantic mappings.
type Theme struct {
	// Base palette
	base     lipgloss.Color
	mantle   lipgloss.Color
	crust    lipgloss.Color
	text     lipgloss.Color
	subtext0 lipgloss.Color
	subtext1 lipgloss.Color
	surface0 lipgloss.Color
	surface1 lipgloss.Color
	surface2 lipgloss.Color
	overlay0 lipgloss.Color
	overlay1 lipgloss.Color
	blue     lipgloss.Color
	lavender lipgloss.Color
	sapphire lipgloss.Color
	sky      lipgloss.Color
	teal     lipgloss.Color
	green    lipgloss.Color
	yellow   lipgloss.Color
	peach    lipgloss.Color
	maroon   lipgloss.Color
	red      lipgloss.Color
	mauve    lipgloss.Color
	pink     lipgloss.Color
	flamingo lipgloss.Color
	rosewater lipgloss.Color

	// Semantic
	correct  lipgloss.Color
	wrong    lipgloss.Color
	pending  lipgloss.Color
	cursor   lipgloss.Color
	title    lipgloss.Color
	wpm      lipgloss.Color
	acc      lipgloss.Color
	timer    lipgloss.Color
	subtle   lipgloss.Color
	hint     lipgloss.Color
	pbFg     lipgloss.Color
	pbBg     lipgloss.Color
	menuSel  lipgloss.Color
	menuSelB lipgloss.Color
	menuDim  lipgloss.Color
	border   lipgloss.Color

	// Syntax
	hlKw  lipgloss.Color
	hlBi  lipgloss.Color
	hlStr lipgloss.Color
	hlCmt lipgloss.Color
	hlNum lipgloss.Color
	hlPun lipgloss.Color

	// Graph
	spark     lipgloss.Color
	sparkPeak lipgloss.Color

	// Heatmap intensity stops (low → high mistake count)
	heat0 lipgloss.Color // no mistakes  — surface
	heat1 lipgloss.Color // low          — yellow
	heat2 lipgloss.Color // mid          — peach
	heat3 lipgloss.Color // high         — maroon
	heat4 lipgloss.Color // critical     — red
	heatBg1 lipgloss.Color
	heatBg2 lipgloss.Color
	heatBg3 lipgloss.Color
	heatBg4 lipgloss.Color
}

// ── Catppuccin Mocha (dark) ────────────────────────────────────────────────

var mocha = Theme{
	base: "#1e1e2e", mantle: "#181825", crust: "#11111b",
	text: "#cdd6f4", subtext0: "#a6adc8", subtext1: "#bac2de",
	surface0: "#313244", surface1: "#45475a", surface2: "#585b70",
	overlay0: "#6c7086", overlay1: "#7f849c",
	blue: "#89b4fa", lavender: "#b4befe", sapphire: "#74c7ec",
	sky: "#89dceb", teal: "#94e2d5", green: "#a6e3a1",
	yellow: "#f9e2af", peach: "#fab387", maroon: "#eba0ac",
	red: "#f38ba8", mauve: "#cba6f7", pink: "#f5c2e7",
	flamingo: "#f2cdcd", rosewater: "#f5e0dc",

	correct: "#a6e3a1", wrong: "#f38ba8", pending: "#45475a",
	cursor: "#cba6f7", title: "#cba6f7",
	wpm: "#f5c2e7", acc: "#94e2d5", timer: "#f9e2af",
	subtle: "#585b70", hint: "#45475a",
	pbFg: "#1e1e2e", pbBg: "#f9e2af",
	menuSel: "#1e1e2e", menuSelB: "#cba6f7", menuDim: "#45475a",
	border: "#313244",
	hlKw: "#cba6f7", hlBi: "#89dceb", hlStr: "#a6e3a1",
	hlCmt: "#6c7086", hlNum: "#fab387", hlPun: "#89b4fa",
	spark: "#cba6f7", sparkPeak: "#f9e2af",
	heat0: "#45475a", heat1: "#f9e2af", heat2: "#fab387",
	heat3: "#eba0ac", heat4: "#f38ba8",
	heatBg1: "#3d3521", heatBg2: "#3d2b15", heatBg3: "#3d1e22", heatBg4: "#3d1515",
}

// ── Catppuccin Latte (light) ───────────────────────────────────────────────

var latte = Theme{
	base: "#eff1f5", mantle: "#e6e9ef", crust: "#dce0e8",
	text: "#4c4f69", subtext0: "#6c6f85", subtext1: "#5c5f77",
	surface0: "#ccd0da", surface1: "#bcc0cc", surface2: "#acb0be",
	overlay0: "#9ca0b0", overlay1: "#8c8fa1",
	blue: "#1e66f5", lavender: "#7287fd", sapphire: "#209fb5",
	sky: "#04a5e5", teal: "#179299", green: "#40a02b",
	yellow: "#df8e1d", peach: "#fe640b", maroon: "#e64553",
	red: "#d20f39", mauve: "#8839ef", pink: "#ea76cb",
	flamingo: "#dd7878", rosewater: "#dc8a78",

	correct: "#40a02b", wrong: "#d20f39", pending: "#acb0be",
	cursor: "#8839ef", title: "#8839ef",
	wpm: "#ea76cb", acc: "#179299", timer: "#df8e1d",
	subtle: "#9ca0b0", hint: "#bcc0cc",
	pbFg: "#eff1f5", pbBg: "#df8e1d",
	menuSel: "#eff1f5", menuSelB: "#8839ef", menuDim: "#9ca0b0",
	border: "#bcc0cc",
	hlKw: "#8839ef", hlBi: "#04a5e5", hlStr: "#40a02b",
	hlCmt: "#9ca0b0", hlNum: "#fe640b", hlPun: "#1e66f5",
	spark: "#8839ef", sparkPeak: "#df8e1d",
	heat0: "#acb0be", heat1: "#df8e1d", heat2: "#fe640b",
	heat3: "#e64553", heat4: "#d20f39",
	heatBg1: "#fdf6e3", heatBg2: "#fdebd0", heatBg3: "#fde2e4", heatBg4: "#fdd6d8",
}

// ── Runtime theme application ──────────────────────────────────────────────

var activeTheme = mocha

func applyTheme(t Theme) {
	activeTheme = t
	correctStyle     = lipgloss.NewStyle().Foreground(t.correct)
	incorrectStyle   = lipgloss.NewStyle().Foreground(t.wrong).Background(t.heatBg4)
	pendingStyle     = lipgloss.NewStyle().Foreground(t.pending)
	// Solid block cursor — bg on the cursor char rather than underline
	cursorStyle      = lipgloss.NewStyle().Foreground(t.base).Background(t.cursor).Bold(true)
	titleStyle       = lipgloss.NewStyle().Foreground(t.mauve).Bold(true)
	wpmStyle         = lipgloss.NewStyle().Foreground(t.wpm).Bold(true)
	accStyle         = lipgloss.NewStyle().Foreground(t.acc).Bold(true)
	timeStyle        = lipgloss.NewStyle().Foreground(t.timer).Bold(true)
	subtleStyle      = lipgloss.NewStyle().Foreground(t.subtext0)
	hintStyle        = lipgloss.NewStyle().Foreground(t.overlay0)
	pbStyle          = lipgloss.NewStyle().Foreground(t.pbFg).Background(t.pbBg).Bold(true)
	errorStyle       = lipgloss.NewStyle().Foreground(t.wrong)
	selectedStyle    = lipgloss.NewStyle().Foreground(t.base).Background(t.mauve).Bold(true).Padding(0, 2).MarginRight(1)
	dimSelectedStyle = lipgloss.NewStyle().Foreground(t.subtext0).Background(t.surface0).Padding(0, 2).MarginRight(1)
	optionStyle      = lipgloss.NewStyle().Foreground(t.overlay0).Padding(0, 2).MarginRight(1)
	hlKeyword        = lipgloss.NewStyle().Foreground(t.hlKw)
	hlBuiltin        = lipgloss.NewStyle().Foreground(t.hlBi)
	hlString         = lipgloss.NewStyle().Foreground(t.hlStr)
	hlComment        = lipgloss.NewStyle().Foreground(t.hlCmt)
	hlNumber         = lipgloss.NewStyle().Foreground(t.hlNum)
	hlPunct          = lipgloss.NewStyle().Foreground(t.hlPun)
	sparkBarStyle    = lipgloss.NewStyle().Foreground(t.spark)
	sparkPeakStyle   = lipgloss.NewStyle().Foreground(t.sparkPeak).Bold(true)
	cardStyle        = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.surface0).
		Padding(1, 3)
}

func init() { applyTheme(mocha) }

func isDark() bool { return activeTheme.base == mocha.base }

// keyHeatStyle returns the correct lipgloss.Style for a key given its mistake count.
func keyHeatStyle(count, maxCount int) lipgloss.Style {
	base := lipgloss.NewStyle().Padding(0, 1).Bold(false)
	if count == 0 || maxCount == 0 {
		return base.Foreground(activeTheme.surface2).Background(activeTheme.surface0)
	}
	h := float64(count) / float64(maxCount)
	switch {
	case h < 0.25:
		return base.Foreground(activeTheme.heat1).Background(activeTheme.heatBg1)
	case h < 0.5:
		return base.Foreground(activeTheme.heat2).Background(activeTheme.heatBg2)
	case h < 0.75:
		return base.Foreground(activeTheme.heat3).Background(activeTheme.heatBg3).Bold(true)
	default:
		return base.Foreground(activeTheme.heat4).Background(activeTheme.heatBg4).Bold(true)
	}
}

// Style vars — all set by applyTheme.
var (
	correctStyle, incorrectStyle, pendingStyle, cursorStyle lipgloss.Style
	titleStyle, wpmStyle, accStyle, timeStyle               lipgloss.Style
	subtleStyle, hintStyle, pbStyle, errorStyle             lipgloss.Style
	selectedStyle, dimSelectedStyle, optionStyle            lipgloss.Style
	hlKeyword, hlBuiltin, hlString, hlComment, hlNumber, hlPunct lipgloss.Style
	sparkBarStyle, sparkPeakStyle                           lipgloss.Style
	cardStyle                                               lipgloss.Style
)

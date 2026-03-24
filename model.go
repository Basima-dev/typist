package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ── Enums ─────────────────────────────────────────────────────────────────────

type appState int

const (
	stateMenu    appState = iota
	stateTyping
	stateResults
	stateHistory
)

type testMode int

const (
	modeWords  testMode = iota
	modeTime
	modeQuote
	modeCode
)

var modeNames = []string{"words", "time", "quote", "code"}
var modeCount = len(modeNames)

// ── Constants ─────────────────────────────────────────────────────────────────

const (
	numWords  = 30
	lineWidth = 60
	visLines  = 3
	histPageSize = 12
)

var timeLimits = []int{15, 30, 60, 120}

// ── Tick / export messages ────────────────────────────────────────────────────

type tickMsg    time.Time
type exportMsg  struct{ path string; err error; format string }

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func exportJSONCmd() tea.Cmd {
	return func() tea.Msg {
		path, err := exportJSON()
		return exportMsg{path: path, err: err, format: "json"}
	}
}

func exportCSVCmd() tea.Cmd {
	return func() tea.Msg {
		path, err := exportCSV()
		return exportMsg{path: path, err: err, format: "csv"}
	}
}

// ── Model ─────────────────────────────────────────────────────────────────────

type Model struct {
	state appState
	mode  testMode

	// time-mode settings
	timeLimitIdx int
	timeLeft     int

	// code-mode settings
	langIdx       int
	activeSnippet Snippet

	// content
	target      []rune
	input       []rune
	activeQuote Quote
	startTime   time.Time
	elapsed     time.Duration
	started     bool

	// features
	blindMode bool

	// stats
	totalKeys int
	errors    int

	// frozen results
	finalWPM float64
	finalAcc float64
	isPB     bool

	// export feedback
	exportMsg string

	// menu cursor (row 0 = mode, row 1 = sub-option)
	menuRow int
	menuCol int

	// history
	histOffset int
	histData   []ScoreEntry

	// layout
	width  int
	height int

	// line wrapping
	lines   []string
	offsets []int
}

func NewModel() Model {
	return Model{
		state:        stateMenu,
		mode:         modeWords,
		timeLimitIdx: 1,
		langIdx:      0,
	}
}

func (m *Model) loadText() {
	var text string
	switch m.mode {
	case modeWords:
		text = generateWords(numWords)
	case modeTime:
		text = generateWords(200)
	case modeQuote:
		m.activeQuote = randomQuote()
		text = m.activeQuote.Text
	case modeCode:
		m.activeSnippet = randomSnippet(langKeys[m.langIdx])
		text = m.activeSnippet.Code
	}
	m.target = []rune(text)
	m.lines, m.offsets = wrapIntoLines(text, lineWidth)
	m.input = nil
	m.totalKeys = 0
	m.errors = 0
	m.started = false
	m.elapsed = 0
	m.exportMsg = ""
	if m.mode == modeTime {
		m.timeLeft = timeLimits[m.timeLimitIdx]
	}
}

func (m Model) modeKey() string {
	return modeNames[int(m.mode)]
}

func (m Model) langKey() string {
	if m.mode == modeCode {
		return langKeys[m.langIdx]
	}
	return ""
}

// ── Init ──────────────────────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd { return nil }

// ── Update ────────────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		if m.state == stateTyping && m.mode == modeTime && m.started {
			m.timeLeft--
			if m.timeLeft <= 0 {
				return m.finishTest(), nil
			}
			return m, tickCmd()
		}

	case exportMsg:
		if msg.err != nil {
			m.exportMsg = errorStyle.Render("export failed: " + msg.err.Error())
		} else {
			m.exportMsg = accStyle.Render("saved → " + msg.path)
		}
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case stateMenu:
			return m.updateMenu(msg)
		case stateTyping:
			return m.updateTyping(msg)
		case stateResults:
			return m.updateResults(msg)
		case stateHistory:
			return m.updateHistory(msg)
		}
	}
	return m, nil
}

// ── Menu ──────────────────────────────────────────────────────────────────────

func (m Model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	subCount := m.subRowCount()

	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit

	case tea.KeyLeft:
		if m.menuRow == 0 {
			m.menuCol = (m.menuCol + modeCount - 1) % modeCount
			m.mode = testMode(m.menuCol)
		} else {
			m.menuCol = (m.menuCol + subCount - 1) % subCount
			m.applySubRow()
		}

	case tea.KeyRight:
		if m.menuRow == 0 {
			m.menuCol = (m.menuCol + 1) % modeCount
			m.mode = testMode(m.menuCol)
		} else {
			m.menuCol = (m.menuCol + 1) % subCount
			m.applySubRow()
		}

	case tea.KeyUp:
		if m.menuRow == 1 {
			m.menuRow = 0
			m.menuCol = int(m.mode)
		}

	case tea.KeyDown:
		if subCount > 0 && m.menuRow == 0 {
			m.menuRow = 1
			if m.mode == modeTime {
				m.menuCol = m.timeLimitIdx
			} else {
				m.menuCol = m.langIdx
			}
		}

	case tea.KeyEnter:
		m.loadText()
		m.state = stateTyping
		if m.mode == modeTime {
			return m, tickCmd()
		}
		return m, nil

	default:
		if len(msg.Runes) > 0 {
			m.loadText()
			m.state = stateTyping
			m.started = true
			nm, cmd := m.handleTypingKey(msg)
			var cmds []tea.Cmd
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			if nm.(Model).mode == modeTime {
				cmds = append(cmds, tickCmd())
			}
			return nm, tea.Batch(cmds...)
		}
	}
	return m, nil
}

func (m Model) subRowCount() int {
	switch m.mode {
	case modeTime:
		return len(timeLimits)
	case modeCode:
		return len(langKeys)
	}
	return 0
}

func (m *Model) applySubRow() {
	switch m.mode {
	case modeTime:
		m.timeLimitIdx = m.menuCol
	case modeCode:
		m.langIdx = m.menuCol
	}
}

// ── Typing ────────────────────────────────────────────────────────────────────

func (m Model) updateTyping(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit
	case tea.KeyCtrlR:
		m.loadText()
		m.state = stateTyping
		if m.mode == modeTime {
			return m, tickCmd()
		}
		return m, nil
	case tea.KeyCtrlB:
		m.blindMode = !m.blindMode
		return m, nil
	}
	return m.handleTypingKey(msg)
}

func (m Model) handleTypingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if !m.started {
		m.started = true
		m.startTime = time.Now()
	}

	switch msg.Type {
	case tea.KeyBackspace:
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}

	case tea.KeySpace:
		m = m.appendRune(' ')

	case tea.KeyTab:
		// In code mode, tab inserts a real tab character
		m = m.appendRune('\t')

	case tea.KeyEnter:
		// In code mode, enter inserts a newline
		if m.mode == modeCode {
			m = m.appendRune('\n')
		}

	case tea.KeyRunes:
		for _, r := range msg.Runes {
			m = m.appendRune(r)
		}
	}

	if m.mode != modeTime && len(m.input) >= len(m.target) {
		return m.finishTest(), nil
	}
	return m, nil
}

func (m Model) appendRune(r rune) Model {
	pos := len(m.input)
	if pos >= len(m.target) {
		return m
	}
	m.totalKeys++
	if r != m.target[pos] {
		m.errors++
	}
	m.input = append(m.input, r)
	return m
}

func (m Model) finishTest() Model {
	m.elapsed = time.Since(m.startTime)
	m.finalWPM = m.calcWPM()
	m.finalAcc = m.calcAccuracy()

	dur := 0
	if m.mode == modeTime {
		dur = timeLimits[m.timeLimitIdx]
	}
	pb := personalBest(m.modeKey(), m.langKey(), dur)
	m.isPB = m.finalWPM > pb

	saveScore(ScoreEntry{
		WPM:      m.finalWPM,
		Accuracy: m.finalAcc,
		Mode:     m.modeKey(),
		Lang:     m.langKey(),
		Duration: dur,
		At:       time.Now(),
	})

	m.state = stateResults
	return m
}

// ── Results ───────────────────────────────────────────────────────────────────

func (m Model) updateResults(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit
	case tea.KeyEnter:
		return m.restart()
	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "r", "R":
			return m.restart()
		case "m", "M":
			next := NewModel()
			next.width = m.width
			next.height = m.height
			next.mode = m.mode
			next.timeLimitIdx = m.timeLimitIdx
			next.langIdx = m.langIdx
			return next, nil
		case "h", "H":
			m.histData = recentSessions(200)
			m.histOffset = 0
			m.state = stateHistory
			return m, nil
		case "j", "J":
			return m, exportJSONCmd()
		case "c", "C":
			return m, exportCSVCmd()
		}
	}
	return m, nil
}

func (m Model) restart() (tea.Model, tea.Cmd) {
	m.loadText()
	m.state = stateTyping
	if m.mode == modeTime {
		return m, tickCmd()
	}
	return m, nil
}

// ── History ───────────────────────────────────────────────────────────────────

func (m Model) updateHistory(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	max := len(m.histData) - histPageSize
	if max < 0 {
		max = 0
	}
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		m.state = stateMenu
		return m, nil
	case tea.KeyDown:
		if m.histOffset < max {
			m.histOffset++
		}
	case tea.KeyUp:
		if m.histOffset > 0 {
			m.histOffset--
		}
	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "q", "Q":
			m.state = stateMenu
			return m, nil
		case "j", "J":
			if m.histOffset < max {
				m.histOffset++
			}
		case "k", "K":
			if m.histOffset > 0 {
				m.histOffset--
			}
		}
	}
	return m, nil
}

// ── Metrics ───────────────────────────────────────────────────────────────────

func (m Model) calcWPM() float64 {
	elapsed := m.elapsed
	if m.state != stateResults {
		if !m.started {
			return 0
		}
		elapsed = time.Since(m.startTime)
	}
	mins := elapsed.Minutes()
	if mins == 0 {
		return 0
	}
	correct := 0
	for i, r := range m.input {
		if i < len(m.target) && r == m.target[i] {
			correct++
		}
	}
	return float64(correct) / 5.0 / mins
}

func (m Model) calcAccuracy() float64 {
	if m.totalKeys == 0 {
		return 100
	}
	return float64(m.totalKeys-m.errors) / float64(m.totalKeys) * 100
}

// ── Views ─────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	switch m.state {
	case stateMenu:
		return m.viewMenu()
	case stateTyping:
		return m.viewTyping()
	case stateResults:
		return m.viewResults()
	case stateHistory:
		return m.viewHistory()
	}
	return ""
}

// ── Menu view ─────────────────────────────────────────────────────────────────

func (m Model) viewMenu() string {
	title := titleStyle.Render("typist")
	sub := subtleStyle.Render("offline · open source · no paywall")

	// Mode row
	var modeBtns []string
	for i, label := range modeNames {
		if i == int(m.mode) {
			modeBtns = append(modeBtns, selectedStyle.Render(" "+label+" "))
		} else {
			modeBtns = append(modeBtns, optionStyle.Render(" "+label+" "))
		}
	}
	modeRow := lipgloss.JoinHorizontal(lipgloss.Center, modeBtns...)

	// Sub-row (time limits or language selector)
	var subRow string
	switch m.mode {
	case modeTime:
		var btns []string
		for i, t := range timeLimits {
			label := fmt.Sprintf("%ds", t)
			if i == m.timeLimitIdx {
				s := dimSelectedStyle
				if m.menuRow == 1 {
					s = selectedStyle
				}
				btns = append(btns, s.Render(" "+label+" "))
			} else {
				btns = append(btns, optionStyle.Render(" "+label+" "))
			}
		}
		subRow = "\n" + lipgloss.JoinHorizontal(lipgloss.Center, btns...)
	case modeCode:
		var btns []string
		for i, lang := range langKeys {
			if i == m.langIdx {
				s := dimSelectedStyle
				if m.menuRow == 1 {
					s = selectedStyle
				}
				btns = append(btns, s.Render(" "+lang+" "))
			} else {
				btns = append(btns, optionStyle.Render(" "+lang+" "))
			}
		}
		subRow = "\n" + lipgloss.JoinHorizontal(lipgloss.Center, btns...)
	}

	var hint string
	if m.mode == modeTime || m.mode == modeCode {
		hint = subtleStyle.Render("← → switch · ↑ ↓ row · enter start · esc quit")
	} else {
		hint = subtleStyle.Render("← → switch · enter start · esc quit")
	}

	body := lipgloss.JoinVertical(lipgloss.Center,
		title, sub, "", modeRow+subRow, "", hint,
	)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, body)
}

// ── Typing view ───────────────────────────────────────────────────────────────

func (m Model) viewTyping() string {
	cursorPos := len(m.input)
	cursorLine := 0
	for i, off := range m.offsets {
		if off <= cursorPos {
			cursorLine = i
		}
	}
	startLine := cursorLine
	if startLine+visLines > len(m.lines) {
		startLine = len(m.lines) - visLines
	}
	if startLine < 0 {
		startLine = 0
	}

	var renderedLines []string
	for li := startLine; li < startLine+visLines && li < len(m.lines); li++ {
		lineStart := m.offsets[li]
		lineRunes := []rune(m.lines[li])
		var sb strings.Builder

		for ci, ch := range lineRunes {
			absPos := lineStart + ci
			display := string(ch)
			if ch == '\t' {
				display = "→   " // visible tab indicator
			}

			if absPos < len(m.input) {
				typed := m.input[absPos]
				correct := typed == ch
				if m.blindMode {
					// Blind mode: show a dot, colored correct/wrong
					if correct {
						sb.WriteString(correctStyle.Render("·"))
					} else {
						sb.WriteString(incorrectStyle.Render("·"))
					}
				} else {
					if correct {
						sb.WriteString(correctStyle.Render(display))
					} else {
						d := string(typed)
						if ch == ' ' || ch == '\t' {
							d = "·"
						}
						sb.WriteString(incorrectStyle.Render(d))
					}
				}
			} else if absPos == cursorPos {
				sb.WriteString(cursorStyle.Render(display))
			} else {
				sb.WriteString(pendingStyle.Render(display))
			}
		}
		renderedLines = append(renderedLines, sb.String())
	}

	textBlock := strings.Join(renderedLines, "\n")

	// Stats bar
	wpmVal := fmt.Sprintf("%.0f", m.calcWPM())
	accVal := fmt.Sprintf("%.0f%%", m.calcAccuracy())

	var timerPart string
	if m.mode == modeTime {
		col := timeStyle
		if m.timeLeft <= 10 {
			col = incorrectStyle
		}
		timerPart = "   " + col.Render(fmt.Sprintf("%ds", m.timeLeft))
	}

	var blindTag string
	if m.blindMode {
		blindTag = "   " + pbStyle.Render(" blind ")
	}

	var langTag string
	if m.mode == modeCode {
		langTag = "   " + subtleStyle.Render(langKeys[m.langIdx])
	}

	stats := lipgloss.JoinHorizontal(lipgloss.Top,
		wpmStyle.Render(wpmVal),
		subtleStyle.Render(" wpm   "),
		accStyle.Render(accVal),
		subtleStyle.Render(" acc"),
		timerPart,
		langTag,
		blindTag,
	)

	var meta string
	switch m.mode {
	case modeQuote:
		meta = subtleStyle.Render("— " + m.activeQuote.Author)
	case modeCode:
		meta = subtleStyle.Render(langKeys[m.langIdx] + " snippet · tab and enter are live")
	}

	hint := hintStyle.Render("ctrl+r restart · ctrl+b blind · esc quit")

	var parts []string
	parts = append(parts, stats, "")
	if meta != "" {
		parts = append(parts, meta)
	}
	parts = append(parts, textBlock, "", hint)

	body := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, body)
}

// ── Results view ──────────────────────────────────────────────────────────────

func (m Model) viewResults() string {
	title := titleStyle.Render("results")

	pbTag := ""
	if m.isPB {
		pbTag = "  " + pbStyle.Render(" new best! ")
	}

	wpmLine := lipgloss.JoinHorizontal(lipgloss.Top,
		wpmStyle.Render(fmt.Sprintf("%-8.0f", m.finalWPM)),
		subtleStyle.Render("wpm"),
		pbTag,
	)
	accLine := lipgloss.JoinHorizontal(lipgloss.Top,
		accStyle.Render(fmt.Sprintf("%-8.1f", m.finalAcc)),
		subtleStyle.Render("accuracy"),
	)
	timeLine := lipgloss.JoinHorizontal(lipgloss.Top,
		timeStyle.Render(fmt.Sprintf("%-8.1f", m.elapsed.Seconds())),
		subtleStyle.Render("seconds"),
	)

	dur := 0
	if m.mode == modeTime {
		dur = timeLimits[m.timeLimitIdx]
	}
	pb := personalBest(m.modeKey(), m.langKey(), dur)
	pbLine := lipgloss.JoinHorizontal(lipgloss.Top,
		subtleStyle.Render(fmt.Sprintf("%-8.0f", pb)),
		subtleStyle.Render("personal best"),
	)

	card := cardStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		title, "",
		wpmLine, accLine, timeLine, "", pbLine,
	))

	actions := lipgloss.JoinVertical(lipgloss.Left,
		pendingStyle.Render("enter / r  → again"),
		pendingStyle.Render("m          → menu"),
		pendingStyle.Render("h          → history"),
		pendingStyle.Render("j          → export json"),
		pendingStyle.Render("c          → export csv"),
		hintStyle.Render("esc        → quit"),
	)

	var exportLine string
	if m.exportMsg != "" {
		exportLine = "\n" + m.exportMsg
	}

	body := lipgloss.JoinVertical(lipgloss.Center, card, "", actions, exportLine)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, body)
}

// ── History view ──────────────────────────────────────────────────────────────

func (m Model) viewHistory() string {
	title := titleStyle.Render("session history")

	if len(m.histData) == 0 {
		body := lipgloss.JoinVertical(lipgloss.Center,
			title, "",
			subtleStyle.Render("no sessions recorded yet"),
			"",
			hintStyle.Render("esc → back"),
		)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, body)
	}

	// Header
	header := lipgloss.JoinHorizontal(lipgloss.Top,
		subtleStyle.Render(fmt.Sprintf("%-8s", "wpm")),
		subtleStyle.Render(fmt.Sprintf("%-8s", "acc%")),
		subtleStyle.Render(fmt.Sprintf("%-10s", "mode")),
		subtleStyle.Render(fmt.Sprintf("%-8s", "lang")),
		subtleStyle.Render("date"),
	)
	divider := subtleStyle.Render(strings.Repeat("─", 48))

	// Rows
	end := m.histOffset + histPageSize
	if end > len(m.histData) {
		end = len(m.histData)
	}

	var rows []string
	for _, e := range m.histData[m.histOffset:end] {
		modeLabel := e.Mode
		if e.Duration > 0 {
			modeLabel += fmt.Sprintf("/%ds", e.Duration)
		}
		lang := e.Lang
		if lang == "" {
			lang = "—"
		}
		row := lipgloss.JoinHorizontal(lipgloss.Top,
			wpmStyle.Render(fmt.Sprintf("%-8.0f", e.WPM)),
			accStyle.Render(fmt.Sprintf("%-8.1f", e.Accuracy)),
			pendingStyle.Render(fmt.Sprintf("%-10s", modeLabel)),
			timeStyle.Render(fmt.Sprintf("%-8s", lang)),
			hintStyle.Render(e.At.Format("Jan 02 15:04")),
		)
		rows = append(rows, row)
	}

	scroll := fmt.Sprintf("%d–%d of %d", m.histOffset+1, end, len(m.histData))
	nav := subtleStyle.Render(scroll + "   j/↓ k/↑ scroll · esc back")

	parts := []string{title, "", header, divider}
	parts = append(parts, rows...)
	parts = append(parts, "", nav)

	body := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, body)
}

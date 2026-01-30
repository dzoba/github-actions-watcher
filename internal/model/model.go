package model

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dzoba/github-actions-watcher/internal/gh"
	"github.com/dzoba/github-actions-watcher/internal/types"
	"github.com/dzoba/github-actions-watcher/internal/ui"
)

// Messages
type repoDetectedMsg struct{ repo string }
type repoErrorMsg struct{ err error }
type runsMsg struct {
	runs []types.WorkflowRun
	json string
}
type runsErrMsg struct{ err error }
type detailMsg struct {
	detail *types.RunDetail
	json   string
}
type detailErrMsg struct{ err error }
type pollTickMsg struct{}
type countdownTickMsg struct{}

// Model is the root Bubbletea model.
type Model struct {
	// Config
	interval time.Duration

	// State
	view              types.View
	repo              string
	repoLoading       bool
	repoError         string
	runs              []types.WorkflowRun
	runsJSON          string
	runsLoading       bool
	runsError         string
	selectedIndex     int
	selectedRunID     int
	detail            *types.RunDetail
	detailJSON        string
	detailLoading     bool
	detailError       string
	detailScrollOffset int
	countdown         int

	// Sub-models
	repoInput textinput.Model

	// Terminal size
	width  int
	height int
}

// New creates a new Model.
func New(interval time.Duration) Model {
	ti := textinput.New()
	ti.Placeholder = "owner/repo"
	ti.CharLimit = 100

	return Model{
		interval:    interval,
		view:        types.ViewList,
		repoLoading: true,
		countdown:   int(interval.Seconds()),
		repoInput:   ti,
	}
}

func (m Model) Init() tea.Cmd {
	return detectRepo
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case repoDetectedMsg:
		m.repo = msg.repo
		m.repoLoading = false
		m.runsLoading = true
		return m, tea.Batch(fetchRuns(m.repo), pollTick(m.interval), countdownTick())

	case repoErrorMsg:
		m.repoLoading = false
		m.repoError = msg.err.Error()
		m.view = types.ViewWelcome
		m.repoInput.Focus()
		return m, m.repoInput.Cursor.BlinkCmd()

	case runsMsg:
		m.runsLoading = false
		m.runsError = ""
		if msg.json != m.runsJSON {
			m.runsJSON = msg.json
			m.runs = msg.runs
		}
		return m, nil

	case runsErrMsg:
		m.runsLoading = false
		m.runsError = msg.err.Error()
		return m, nil

	case detailMsg:
		m.detailLoading = false
		m.detailError = ""
		if msg.json != m.detailJSON {
			m.detailJSON = msg.json
			m.detail = msg.detail
		}
		return m, nil

	case detailErrMsg:
		m.detailLoading = false
		m.detailError = msg.err.Error()
		return m, nil

	case pollTickMsg:
		m.countdown = int(m.interval.Seconds())
		cmds := []tea.Cmd{pollTick(m.interval)}
		if m.repo != "" {
			cmds = append(cmds, fetchRuns(m.repo))
			if m.view == types.ViewDetail && m.selectedRunID != 0 {
				cmds = append(cmds, fetchRunDetail(m.repo, m.selectedRunID))
			}
		}
		return m, tea.Batch(cmds...)

	case countdownTickMsg:
		if m.countdown > 0 {
			m.countdown--
		}
		return m, countdownTick()

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// Pass through to text input when in repo-input or welcome view
	if m.view == types.ViewRepoInput || m.view == types.ViewWelcome {
		var cmd tea.Cmd
		m.repoInput, cmd = m.repoInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.view {
	case types.ViewList:
		return m.handleListKey(msg)
	case types.ViewDetail:
		return m.handleDetailKey(msg)
	case types.ViewRepoInput:
		return m.handleRepoInputKey(msg)
	case types.ViewWelcome:
		return m.handleWelcomeKey(msg)
	}
	return m, nil
}

func (m Model) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, ui.ListKeys.Quit):
		return m, tea.Quit
	case key.Matches(msg, ui.ListKeys.Up):
		if m.selectedIndex > 0 {
			m.selectedIndex--
		}
	case key.Matches(msg, ui.ListKeys.Down):
		if m.selectedIndex < len(m.runs)-1 {
			m.selectedIndex++
		}
	case key.Matches(msg, ui.ListKeys.Enter):
		if len(m.runs) > 0 && m.selectedIndex < len(m.runs) {
			run := m.runs[m.selectedIndex]
			m.selectedRunID = run.DatabaseID
			m.detail = nil
			m.detailJSON = ""
			m.detailScrollOffset = 0
			m.detailLoading = true
			m.view = types.ViewDetail
			return m, fetchRunDetail(m.repo, run.DatabaseID)
		}
	case key.Matches(msg, ui.ListKeys.Switch):
		m.view = types.ViewRepoInput
		m.repoInput.SetValue(m.repo)
		m.repoInput.Focus()
		return m, m.repoInput.Cursor.BlinkCmd()
	case key.Matches(msg, ui.ListKeys.Refresh):
		m.countdown = int(m.interval.Seconds())
		return m, fetchRuns(m.repo)
	}
	return m, nil
}

func (m Model) handleDetailKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, ui.DetailKeys.Quit):
		return m, tea.Quit
	case key.Matches(msg, ui.DetailKeys.Back):
		m.view = types.ViewList
		m.selectedRunID = 0
		m.detail = nil
	case key.Matches(msg, ui.DetailKeys.Up):
		if m.detailScrollOffset > 0 {
			m.detailScrollOffset--
		}
	case key.Matches(msg, ui.DetailKeys.Down):
		m.detailScrollOffset++
	case key.Matches(msg, ui.DetailKeys.Open):
		if m.detail != nil && m.detail.URL != "" {
			openBrowser(m.detail.URL)
		}
	case key.Matches(msg, ui.DetailKeys.Refresh):
		m.countdown = int(m.interval.Seconds())
		return m, tea.Batch(fetchRuns(m.repo), fetchRunDetail(m.repo, m.selectedRunID))
	}
	return m, nil
}

func (m Model) handleRepoInputKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, ui.RepoInputKeys.Cancel):
		m.view = types.ViewList
		m.repoInput.Blur()
		return m, nil
	case key.Matches(msg, ui.RepoInputKeys.Confirm):
		val := strings.TrimSpace(m.repoInput.Value())
		if val != "" && strings.Contains(val, "/") {
			m.repo = val
			m.runs = nil
			m.runsJSON = ""
			m.selectedIndex = 0
			m.view = types.ViewList
			m.repoInput.Blur()
			m.runsLoading = true
			return m, fetchRuns(m.repo)
		}
		return m, nil
	}
	// Pass to text input
	var cmd tea.Cmd
	m.repoInput, cmd = m.repoInput.Update(msg)
	return m, cmd
}

func (m Model) handleWelcomeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.String() == "q":
		return m, tea.Quit
	case key.Matches(msg, ui.RepoInputKeys.Cancel):
		return m, tea.Quit
	case key.Matches(msg, ui.RepoInputKeys.Confirm):
		val := strings.TrimSpace(m.repoInput.Value())
		if val != "" && strings.Contains(val, "/") {
			m.repo = val
			m.repoError = ""
			m.view = types.ViewList
			m.repoInput.Blur()
			m.runsLoading = true
			return m, tea.Batch(fetchRuns(m.repo), pollTick(m.interval), countdownTick())
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.repoInput, cmd = m.repoInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.repoLoading {
		return ui.Dim.Render("Detecting repository...")
	}

	switch m.view {
	case types.ViewWelcome:
		return m.welcomeView()
	case types.ViewList:
		return m.listViewFull()
	case types.ViewDetail:
		return m.detailViewFull()
	case types.ViewRepoInput:
		return m.repoInputViewFull()
	}
	return ""
}

func (m Model) listViewFull() string {
	var b strings.Builder
	// Header
	b.WriteString(ui.CyanBold.Render("GitHub Actions"))
	b.WriteString(" - ")
	b.WriteString(ui.Bold.Render(m.repo))
	b.WriteString("\n\n")

	if m.runsError != "" {
		b.WriteString(ui.Red.Render("Error: " + m.runsError))
		b.WriteString("\n")
	}

	if m.runsLoading && len(m.runs) == 0 {
		b.WriteString(ui.Dim.Render("Loading runs..."))
	} else {
		b.WriteString(m.listView())
	}

	b.WriteString("\n")
	b.WriteString(m.footerView())
	return b.String()
}

func (m Model) detailViewFull() string {
	var b strings.Builder
	// Header
	b.WriteString(ui.CyanBold.Render("GitHub Actions"))
	b.WriteString(" - ")
	b.WriteString(ui.Bold.Render(m.repo))
	b.WriteString("\n\n")

	b.WriteString(m.detailView())

	b.WriteString("\n")
	b.WriteString(m.footerView())
	return b.String()
}

func (m Model) repoInputViewFull() string {
	var b strings.Builder
	b.WriteString(ui.CyanBold.Render("GitHub Actions"))
	b.WriteString(" - ")
	b.WriteString(ui.Bold.Render(m.repo))
	b.WriteString("\n\n")
	b.WriteString(m.repoInputView())
	b.WriteString("\n")
	b.WriteString(ui.Dim.Render("enter: confirm | esc: cancel"))
	return b.String()
}

func (m Model) footerView() string {
	var hint string
	switch m.view {
	case types.ViewList:
		hint = "up/down: navigate | enter: details | s: switch repo | r: refresh | q: quit"
	case types.ViewDetail:
		hint = "up/down: scroll | esc: back | o: open in browser | r: refresh | q: quit"
	}
	return ui.Dim.Render(fmt.Sprintf("%s | next refresh: %ds", hint, m.countdown))
}

// Commands

func detectRepo() tea.Msg {
	repo, err := gh.DetectRepo()
	if err != nil {
		return repoErrorMsg{err}
	}
	return repoDetectedMsg{repo}
}

func fetchRuns(repo string) tea.Cmd {
	return func() tea.Msg {
		runs, err := gh.FetchRuns(repo)
		if err != nil {
			return runsErrMsg{err}
		}
		j, _ := json.Marshal(runs)
		return runsMsg{runs: runs, json: string(j)}
	}
}

func fetchRunDetail(repo string, runID int) tea.Cmd {
	return func() tea.Msg {
		detail, err := gh.FetchRunDetail(repo, runID)
		if err != nil {
			return detailErrMsg{err}
		}
		j, _ := json.Marshal(detail)
		return detailMsg{detail: detail, json: string(j)}
	}
}

func pollTick(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return pollTickMsg{}
	})
}

func countdownTick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return countdownTickMsg{}
	})
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return
	}
	_ = cmd.Start()
}

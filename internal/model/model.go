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

// repoTab holds all per-repo state.
type repoTab struct {
	repo               string
	runs               []types.WorkflowRun
	runsJSON           string
	runsLoading        bool
	runsError          string
	selectedIndex      int
	selectedRunID      int
	detail             *types.RunDetail
	detailJSON         string
	detailLoading      bool
	detailError        string
	detailScrollOffset int
	view               types.View // ViewList or ViewDetail (per-tab)
}

// Messages
type repoDetectedMsg struct{ repo string }
type repoErrorMsg struct{ err error }
type runsMsg struct {
	tabIndex int
	runs     []types.WorkflowRun
	json     string
}
type runsErrMsg struct {
	tabIndex int
	err      error
}
type detailMsg struct {
	tabIndex int
	detail   *types.RunDetail
	json     string
}
type detailErrMsg struct {
	tabIndex int
	err      error
}
type pollTickMsg struct{}
type countdownTickMsg struct{}
type repoListMsg struct{ repos []types.PickerRepo }
type repoListErrMsg struct{ err error }

// Model is the root Bubbletea model.
type Model struct {
	// Config
	interval time.Duration

	// Tabs
	tabs      []repoTab
	activeTab int

	// Root-level state
	repoLoading bool
	repoError   string
	countdown   int

	// Picker state
	showPicker     bool
	pickerRepos    []types.PickerRepo
	pickerLoading  bool
	pickerSelected int
	pickerFilter   textinput.Model

	// Terminal size
	width  int
	height int
}

// New creates a new Model.
func New(interval time.Duration) Model {
	ti := textinput.New()
	ti.Placeholder = "filter or owner/repo"
	ti.CharLimit = 100

	return Model{
		interval:    interval,
		repoLoading: true,
		countdown:   int(interval.Seconds()),
		pickerFilter: ti,
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
		m.repoLoading = false
		tab := repoTab{
			repo:        msg.repo,
			runsLoading: true,
			view:        types.ViewList,
		}
		m.tabs = []repoTab{tab}
		m.activeTab = 0
		return m, tea.Batch(fetchRuns(msg.repo, 0), pollTick(m.interval), countdownTick())

	case repoErrorMsg:
		m.repoLoading = false
		m.repoError = msg.err.Error()
		m.showPicker = true
		m.pickerLoading = true
		m.pickerFilter.Focus()
		return m, tea.Batch(m.pickerFilter.Cursor.BlinkCmd(), fetchRepoList())

	case runsMsg:
		if msg.tabIndex < len(m.tabs) {
			t := &m.tabs[msg.tabIndex]
			t.runsLoading = false
			t.runsError = ""
			if msg.json != t.runsJSON {
				t.runsJSON = msg.json
				t.runs = msg.runs
			}
		}
		return m, nil

	case runsErrMsg:
		if msg.tabIndex < len(m.tabs) {
			t := &m.tabs[msg.tabIndex]
			t.runsLoading = false
			t.runsError = msg.err.Error()
		}
		return m, nil

	case detailMsg:
		if msg.tabIndex < len(m.tabs) {
			t := &m.tabs[msg.tabIndex]
			t.detailLoading = false
			t.detailError = ""
			if msg.json != t.detailJSON {
				t.detailJSON = msg.json
				t.detail = msg.detail
			}
		}
		return m, nil

	case detailErrMsg:
		if msg.tabIndex < len(m.tabs) {
			t := &m.tabs[msg.tabIndex]
			t.detailLoading = false
			t.detailError = msg.err.Error()
		}
		return m, nil

	case repoListMsg:
		m.pickerLoading = false
		m.pickerRepos = msg.repos
		return m, nil

	case repoListErrMsg:
		m.pickerLoading = false
		return m, nil

	case pollTickMsg:
		m.countdown = int(m.interval.Seconds())
		cmds := []tea.Cmd{pollTick(m.interval)}
		for i := range m.tabs {
			cmds = append(cmds, fetchRuns(m.tabs[i].repo, i))
			if i == m.activeTab && m.tabs[i].view == types.ViewDetail && m.tabs[i].selectedRunID != 0 {
				cmds = append(cmds, fetchRunDetail(m.tabs[i].repo, m.tabs[i].selectedRunID, i))
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

	// Pass through to text input when picker is showing
	if m.showPicker {
		var cmd tea.Cmd
		m.pickerFilter, cmd = m.pickerFilter.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Ctrl+C always quits regardless of view
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	if m.showPicker {
		return m.handlePickerKey(msg)
	}

	// Handle welcome screen (no tabs)
	if len(m.tabs) == 0 {
		return m.handleWelcomeKey(msg)
	}

	// Number keys 1-9 for tab switching
	if len(m.tabs) > 1 {
		k := msg.String()
		if len(k) == 1 && k[0] >= '1' && k[0] <= '9' {
			idx := int(k[0] - '1')
			if idx < len(m.tabs) {
				m.activeTab = idx
				return m, nil
			}
		}
	}

	tab := m.tabs[m.activeTab]
	switch tab.view {
	case types.ViewList:
		return m.handleListKey(msg)
	case types.ViewDetail:
		return m.handleDetailKey(msg)
	}
	return m, nil
}

func (m Model) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	t := &m.tabs[m.activeTab]
	switch {
	case key.Matches(msg, ui.ListKeys.Quit):
		return m, tea.Quit
	case key.Matches(msg, ui.ListKeys.Tab):
		if len(m.tabs) > 1 {
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
		}
	case key.Matches(msg, ui.ListKeys.ShiftTab):
		if len(m.tabs) > 1 {
			m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
		}
	case key.Matches(msg, ui.ListKeys.CloseTab):
		if len(m.tabs) > 1 {
			return m.closeTab(m.activeTab), nil
		}
	case key.Matches(msg, ui.ListKeys.Up):
		if t.selectedIndex > 0 {
			t.selectedIndex--
		}
	case key.Matches(msg, ui.ListKeys.Down):
		if t.selectedIndex < len(t.runs)-1 {
			t.selectedIndex++
		}
	case key.Matches(msg, ui.ListKeys.Enter):
		if len(t.runs) > 0 && t.selectedIndex < len(t.runs) {
			run := t.runs[t.selectedIndex]
			t.selectedRunID = run.DatabaseID
			t.detail = nil
			t.detailJSON = ""
			t.detailScrollOffset = 0
			t.detailLoading = true
			t.view = types.ViewDetail
			return m, fetchRunDetail(t.repo, run.DatabaseID, m.activeTab)
		}
	case key.Matches(msg, ui.ListKeys.Switch):
		m.showPicker = true
		m.pickerLoading = true
		m.pickerSelected = 0
		m.pickerFilter.SetValue("")
		m.pickerFilter.Focus()
		return m, tea.Batch(m.pickerFilter.Cursor.BlinkCmd(), fetchRepoList())
	case key.Matches(msg, ui.ListKeys.Refresh):
		m.countdown = int(m.interval.Seconds())
		return m, fetchRuns(t.repo, m.activeTab)
	}
	return m, nil
}

func (m Model) handleDetailKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	t := &m.tabs[m.activeTab]
	switch {
	case key.Matches(msg, ui.DetailKeys.Quit):
		return m, tea.Quit
	case key.Matches(msg, ui.DetailKeys.Tab):
		if len(m.tabs) > 1 {
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
		}
	case key.Matches(msg, ui.DetailKeys.ShiftTab):
		if len(m.tabs) > 1 {
			m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
		}
	case key.Matches(msg, ui.DetailKeys.CloseTab):
		if len(m.tabs) > 1 {
			return m.closeTab(m.activeTab), nil
		}
	case key.Matches(msg, ui.DetailKeys.Back):
		t.view = types.ViewList
		t.selectedRunID = 0
		t.detail = nil
	case key.Matches(msg, ui.DetailKeys.Up):
		if t.detailScrollOffset > 0 {
			t.detailScrollOffset--
		}
	case key.Matches(msg, ui.DetailKeys.Down):
		t.detailScrollOffset++
	case key.Matches(msg, ui.DetailKeys.Open):
		if t.detail != nil && t.detail.URL != "" {
			openBrowser(t.detail.URL)
		}
	case key.Matches(msg, ui.DetailKeys.Refresh):
		m.countdown = int(m.interval.Seconds())
		return m, tea.Batch(fetchRuns(t.repo, m.activeTab), fetchRunDetail(t.repo, t.selectedRunID, m.activeTab))
	}
	return m, nil
}

func (m Model) handlePickerKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, ui.PickerKeys.Cancel):
		m.showPicker = false
		m.pickerFilter.Blur()
		// If no tabs exist (welcome screen), quit
		if len(m.tabs) == 0 {
			return m, tea.Quit
		}
		return m, nil
	case key.Matches(msg, ui.PickerKeys.Enter):
		return m.pickerSelect()
	case key.Matches(msg, ui.PickerKeys.Remove):
		return m.pickerRemoveTab()
	case key.Matches(msg, ui.PickerKeys.Up):
		if m.pickerSelected > 0 {
			m.pickerSelected--
		}
		return m, nil
	case key.Matches(msg, ui.PickerKeys.Down):
		filtered := m.filteredPickerRepos()
		if m.pickerSelected < len(filtered)-1 {
			m.pickerSelected++
		}
		return m, nil
	}
	// Pass to text input for typing
	var cmd tea.Cmd
	m.pickerFilter, cmd = m.pickerFilter.Update(msg)
	// Reset selection when filter changes
	m.pickerSelected = 0
	return m, cmd
}

func (m Model) handleWelcomeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "q" {
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) closeTab(idx int) Model {
	m.tabs = append(m.tabs[:idx], m.tabs[idx+1:]...)
	if m.activeTab >= len(m.tabs) {
		m.activeTab = len(m.tabs) - 1
	}
	return m
}

func (m Model) filteredPickerRepos() []types.PickerRepo {
	filter := strings.ToLower(strings.TrimSpace(m.pickerFilter.Value()))
	if filter == "" {
		return m.pickerRepos
	}
	var result []types.PickerRepo
	for _, r := range m.pickerRepos {
		if strings.Contains(strings.ToLower(r.NameWithOwner), filter) {
			result = append(result, r)
		}
	}
	return result
}

func (m Model) isRepoOpen(name string) bool {
	for _, t := range m.tabs {
		if t.repo == name {
			return true
		}
	}
	return false
}

func (m Model) tabIndexForRepo(name string) int {
	for i, t := range m.tabs {
		if t.repo == name {
			return i
		}
	}
	return -1
}

func (m Model) pickerSelect() (tea.Model, tea.Cmd) {
	filtered := m.filteredPickerRepos()
	var repoName string

	if len(filtered) > 0 && m.pickerSelected < len(filtered) {
		repoName = filtered[m.pickerSelected].NameWithOwner
	} else {
		// Manual entry: typed text that contains "/"
		val := strings.TrimSpace(m.pickerFilter.Value())
		if val != "" && strings.Contains(val, "/") {
			repoName = val
		} else {
			return m, nil
		}
	}

	// If already open, switch to it
	if idx := m.tabIndexForRepo(repoName); idx >= 0 {
		m.activeTab = idx
		m.showPicker = false
		m.pickerFilter.Blur()
		return m, nil
	}

	// Add new tab
	tab := repoTab{
		repo:        repoName,
		runsLoading: true,
		view:        types.ViewList,
	}
	m.tabs = append(m.tabs, tab)
	newIdx := len(m.tabs) - 1
	m.activeTab = newIdx
	m.showPicker = false
	m.pickerFilter.Blur()

	cmds := []tea.Cmd{fetchRuns(repoName, newIdx)}
	// Start polling if this is the first tab
	if len(m.tabs) == 1 {
		cmds = append(cmds, pollTick(m.interval), countdownTick())
	}
	return m, tea.Batch(cmds...)
}

func (m Model) pickerRemoveTab() (tea.Model, tea.Cmd) {
	filtered := m.filteredPickerRepos()
	if len(filtered) == 0 || m.pickerSelected >= len(filtered) {
		return m, nil
	}
	repoName := filtered[m.pickerSelected].NameWithOwner
	idx := m.tabIndexForRepo(repoName)
	if idx < 0 || len(m.tabs) <= 1 {
		return m, nil
	}
	m = m.closeTab(idx)
	return m, nil
}

func (m Model) View() string {
	if m.repoLoading {
		return ui.Dim.Render("Detecting repository...")
	}

	if m.showPicker {
		return m.pickerViewFull()
	}

	if len(m.tabs) == 0 {
		return m.welcomeView()
	}

	t := m.tabs[m.activeTab]
	switch t.view {
	case types.ViewList:
		return m.listViewFull()
	case types.ViewDetail:
		return m.detailViewFull()
	}
	return ""
}

func (m Model) tabBar() string {
	if len(m.tabs) <= 1 {
		return ""
	}
	var parts []string
	for i, t := range m.tabs {
		label := fmt.Sprintf("%d: %s", i+1, t.repo)
		if i == m.activeTab {
			parts = append(parts, ui.TabActive.Render(label))
		} else {
			parts = append(parts, ui.TabInactive.Render(label))
		}
	}
	return strings.Join(parts, "  ") + "\n"
}

func (m Model) listViewFull() string {
	t := m.tabs[m.activeTab]
	var b strings.Builder

	// Tab bar
	b.WriteString(m.tabBar())

	// Header
	b.WriteString(ui.CyanBold.Render("GitHub Actions"))
	b.WriteString(" - ")
	b.WriteString(ui.Bold.Render(t.repo))
	b.WriteString("\n\n")

	if t.runsError != "" {
		b.WriteString(ui.Red.Render("Error: " + t.runsError))
		b.WriteString("\n")
	}

	if t.runsLoading && len(t.runs) == 0 {
		b.WriteString(ui.Dim.Render("Loading runs..."))
	} else {
		b.WriteString(m.listView())
	}

	b.WriteString("\n")
	b.WriteString(m.footerView())
	return b.String()
}

func (m Model) detailViewFull() string {
	t := m.tabs[m.activeTab]
	var b strings.Builder

	// Tab bar
	b.WriteString(m.tabBar())

	// Header
	b.WriteString(ui.CyanBold.Render("GitHub Actions"))
	b.WriteString(" - ")
	b.WriteString(ui.Bold.Render(t.repo))
	b.WriteString("\n\n")

	b.WriteString(m.detailView())

	b.WriteString("\n")
	b.WriteString(m.footerView())
	return b.String()
}

func (m Model) pickerViewFull() string {
	var b strings.Builder
	b.WriteString(ui.CyanBold.Render("GitHub Actions"))
	b.WriteString(" - ")
	b.WriteString(ui.Bold.Render("Select Repository"))
	b.WriteString("\n\n")
	b.WriteString(m.pickerView())
	return b.String()
}

func (m Model) footerView() string {
	t := m.tabs[m.activeTab]
	var hint string
	switch t.view {
	case types.ViewList:
		hint = "up/down: navigate | enter: details | s: switch repo | r: refresh | q: quit"
		if len(m.tabs) > 1 {
			hint = "up/down: navigate | enter: details | tab/shift-tab: switch tab | w: close tab | s: add repo | r: refresh | q: quit"
		}
	case types.ViewDetail:
		hint = "up/down: scroll | esc: back | o: open in browser | r: refresh | q: quit"
		if len(m.tabs) > 1 {
			hint = "up/down: scroll | esc: back | tab/shift-tab: switch tab | o: open | r: refresh | q: quit"
		}
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

func fetchRuns(repo string, tabIndex int) tea.Cmd {
	return func() tea.Msg {
		runs, err := gh.FetchRuns(repo)
		if err != nil {
			return runsErrMsg{tabIndex: tabIndex, err: err}
		}
		j, _ := json.Marshal(runs)
		return runsMsg{tabIndex: tabIndex, runs: runs, json: string(j)}
	}
}

func fetchRunDetail(repo string, runID int, tabIndex int) tea.Cmd {
	return func() tea.Msg {
		detail, err := gh.FetchRunDetail(repo, runID)
		if err != nil {
			return detailErrMsg{tabIndex: tabIndex, err: err}
		}
		j, _ := json.Marshal(detail)
		return detailMsg{tabIndex: tabIndex, detail: detail, json: string(j)}
	}
}

func fetchRepoList() tea.Cmd {
	return func() tea.Msg {
		repos, err := gh.FetchRepoList()
		if err != nil {
			return repoListErrMsg{err}
		}
		return repoListMsg{repos}
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

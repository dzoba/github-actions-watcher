package model

import (
	"strings"

	"github.com/dzoba/github-actions-watcher/internal/format"
	"github.com/dzoba/github-actions-watcher/internal/ui"
)

func (m Model) pickerView() string {
	var b strings.Builder

	// Filter input
	b.WriteString(ui.Dim.Render("Search: "))
	b.WriteString(m.pickerFilter.View())
	b.WriteString("\n\n")

	if m.pickerLoading {
		b.WriteString(ui.Dim.Render("Loading repositories..."))
		b.WriteString("\n")
	} else {
		filtered := m.filteredPickerRepos()
		if len(filtered) == 0 {
			filterVal := strings.TrimSpace(m.pickerFilter.Value())
			if filterVal != "" && strings.Contains(filterVal, "/") {
				b.WriteString(ui.Dim.Render("No matches. Press enter to add "))
				b.WriteString(ui.Bold.Render(filterVal))
				b.WriteString("\n")
			} else {
				b.WriteString(ui.Dim.Render("No repositories found."))
				b.WriteString("\n")
			}
		} else {
			for i, repo := range filtered {
				// Selector
				if i == m.pickerSelected {
					b.WriteString("> ")
				} else {
					b.WriteString("  ")
				}

				// Repo name
				name := repo.NameWithOwner
				if m.isRepoOpen(name) {
					b.WriteString(ui.CyanBold.Render(name))
					b.WriteString(ui.Dim.Render(" *"))
				} else {
					b.WriteString(name)
				}

				// Pushed-at time
				if repo.PushedAt != "" {
					b.WriteString("  ")
					b.WriteString(ui.Dim.Render(format.RelativeTime(repo.PushedAt)))
				}

				b.WriteByte('\n')
			}
		}
	}

	b.WriteString("\n")
	b.WriteString(ui.Dim.Render("up/down: navigate | enter: select | x: remove tab | esc: cancel"))
	return b.String()
}

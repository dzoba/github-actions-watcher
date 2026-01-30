package model

import (
	"strings"

	"github.com/dzoba/github-actions-watcher/internal/ui"
)

func (m Model) welcomeView() string {
	var b strings.Builder
	b.WriteString(ui.CyanBold.Render("GitHub Actions Watcher"))
	b.WriteString("\n\n")
	b.WriteString(ui.Yellow.Render("No GitHub repository detected in this directory."))
	b.WriteString("\n")
	b.WriteString(ui.Dim.Render("To get started, either:"))
	b.WriteString("\n")
	b.WriteString("  cd into a directory with a GitHub remote and run ")
	b.WriteString(ui.Bold.Render("ghaw"))
	b.WriteString("\n")
	b.WriteString("  or press any key to open the repo picker")
	b.WriteString("\n\n")
	b.WriteString(ui.Dim.Render("Loading repo picker..."))
	b.WriteString("\n\n")
	b.WriteString(ui.Dim.Render("q: quit"))
	return b.String()
}

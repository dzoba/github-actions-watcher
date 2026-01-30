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
	b.WriteString("  or type a repo below (e.g. ")
	b.WriteString(ui.Bold.Render("owner/repo"))
	b.WriteString(")")
	b.WriteString("\n\n")
	b.WriteString(ui.Dim.Render("owner/repo: "))
	b.WriteString(m.repoInput.View())
	b.WriteString("\n\n")
	b.WriteString(ui.Dim.Render("enter: confirm | esc/q: quit"))
	return b.String()
}

package model

import (
	"strings"

	"github.com/dzoba/github-actions-watcher/internal/ui"
)

func (m Model) repoInputView() string {
	var b strings.Builder
	b.WriteString("Switch repository:\n")
	b.WriteString(ui.Dim.Render("owner/repo: "))
	b.WriteString(m.repoInput.View())
	return b.String()
}

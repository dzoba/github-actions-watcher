package model

import (
	"strings"

	"github.com/dzoba/github-actions-watcher/internal/format"
	"github.com/dzoba/github-actions-watcher/internal/ui"
)

func (m Model) listView() string {
	if len(m.runs) == 0 {
		return ui.Dim.Render("No workflow runs found.")
	}

	cols := m.width
	if cols == 0 {
		cols = 120
	}

	showTime := cols >= 70
	showBranch := cols >= 55
	showWorkflow := cols >= 40

	titleMax := cols - (2 + 14)
	if showWorkflow {
		titleMax -= 21
	}
	if showBranch {
		titleMax -= 19
	}
	if showTime {
		titleMax -= 9
	}
	if titleMax < 15 {
		titleMax = 15
	}

	var b strings.Builder
	for i, run := range m.runs {
		if i > 0 {
			b.WriteByte('\n')
		}

		// Selector
		if i == m.selectedIndex {
			b.WriteString("> ")
		} else {
			b.WriteString("  ")
		}

		// Status badge (pad plain text to 14, then style)
		badgeText, badgeColor := format.StatusBadge(run.Status, run.Conclusion)
		b.WriteString(ui.BadgeStyle(badgeColor).Render(format.Pad(badgeText, 13)))
		b.WriteByte(' ')

		// Workflow name (pad plain text, then style)
		if showWorkflow {
			b.WriteString(ui.Blue.Render(format.Pad(format.Truncate(run.WorkflowName, 20), 20)))
			b.WriteByte(' ')
		}

		// Branch (pad plain text, then style)
		if showBranch {
			b.WriteString(ui.Magenta.Render(format.Pad(format.Truncate(run.HeadBranch, 18), 18)))
			b.WriteByte(' ')
		}

		// Title (plain text)
		b.WriteString(format.Pad(format.Truncate(run.DisplayTitle, titleMax), titleMax))

		// Time (pad plain text, then style)
		if showTime {
			b.WriteByte(' ')
			b.WriteString(ui.Dim.Render(format.Pad(format.RelativeTime(run.CreatedAt), 8)))
		}
	}
	return b.String()
}

package model

import (
	"fmt"
	"strings"

	"github.com/dzoba/github-actions-watcher/internal/format"
	"github.com/dzoba/github-actions-watcher/internal/ui"
)

func (m Model) detailView() string {
	t := m.tabs[m.activeTab]
	if t.detailLoading && t.detail == nil {
		return ui.Dim.Render("Loading run details...")
	}
	if t.detailError != "" {
		return ui.Red.Render("Error: " + t.detailError)
	}
	if t.detail == nil {
		return ui.Dim.Render("No detail available.")
	}

	d := t.detail
	var lines []string

	// Title line
	badgeText, badgeColor := format.StatusBadge(d.Status, d.Conclusion)
	titleLine := ui.BadgeStyle(badgeColor).Render(badgeText) + " " + ui.Bold.Render(d.DisplayTitle)
	lines = append(lines, titleLine)

	// Metadata line
	meta := fmt.Sprintf("%s #%d on %s (%s)", d.WorkflowName, d.Number, d.HeadBranch, d.Event)
	if t.detailLoading {
		meta += " fetching..."
	}
	lines = append(lines, ui.Dim.Render(meta))

	// Separator
	lines = append(lines, ui.Dim.Render("---"))

	// Jobs and steps
	for _, job := range d.Jobs {
		badgeText, badgeColor := format.StatusBadge(job.Status, job.Conclusion)
		jobLine := ui.BadgeStyle(badgeColor).Render(badgeText) + " " + ui.Bold.Render(job.Name)
		if job.StartedAt != "" {
			jobLine += " " + ui.Dim.Render("("+format.Duration(job.StartedAt, job.CompletedAt)+")")
		}
		lines = append(lines, "")
		lines = append(lines, jobLine)

		for _, step := range job.Steps {
			sbText, sbColor := format.StatusBadge(step.Status, step.Conclusion)
			stepLine := "  " + ui.BadgeStyle(sbColor).Render(sbText) + " " + step.Name
			if step.StartedAt != "" {
				stepLine += " " + ui.Dim.Render("("+format.Duration(step.StartedAt, step.CompletedAt)+")")
			}
			lines = append(lines, stepLine)
		}
	}

	// Apply scroll offset
	if t.detailScrollOffset > 0 && t.detailScrollOffset < len(lines) {
		lines = lines[t.detailScrollOffset:]
	} else if t.detailScrollOffset >= len(lines) {
		lines = nil
	}

	return strings.Join(lines, "\n")
}

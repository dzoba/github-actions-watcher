package model

import (
	"fmt"
	"strings"

	"github.com/dzoba/github-actions-watcher/internal/format"
	"github.com/dzoba/github-actions-watcher/internal/ui"
)

func (m Model) detailView() string {
	if m.detailLoading && m.detail == nil {
		return ui.Dim.Render("Loading run details...")
	}
	if m.detailError != "" {
		return ui.Red.Render("Error: " + m.detailError)
	}
	if m.detail == nil {
		return ui.Dim.Render("No detail available.")
	}

	d := m.detail
	var lines []string

	// Title line
	badgeText, badgeColor := format.StatusBadge(d.Status, d.Conclusion)
	titleLine := ui.BadgeStyle(badgeColor).Render(badgeText) + " " + ui.Bold.Render(d.DisplayTitle)
	lines = append(lines, titleLine)

	// Metadata line
	meta := fmt.Sprintf("%s #%d on %s (%s)", d.WorkflowName, d.Number, d.HeadBranch, d.Event)
	if m.detailLoading {
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
			if step.StartedAt != "" && step.CompletedAt != "" {
				stepLine += " " + ui.Dim.Render("("+format.Duration(step.StartedAt, step.CompletedAt)+")")
			}
			lines = append(lines, stepLine)
		}
	}

	// Apply scroll offset
	if m.detailScrollOffset > 0 && m.detailScrollOffset < len(lines) {
		lines = lines[m.detailScrollOffset:]
	} else if m.detailScrollOffset >= len(lines) {
		lines = nil
	}

	return strings.Join(lines, "\n")
}

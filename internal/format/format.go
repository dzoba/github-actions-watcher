package format

import (
	"fmt"
	"math"
	"time"

	"github.com/dzoba/github-actions-watcher/internal/types"
)

// RelativeTime returns a human-friendly relative time string (e.g. "5m ago").
func RelativeTime(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return ""
	}
	diff := time.Since(t)
	sec := int(math.Floor(diff.Seconds()))
	if sec < 0 {
		sec = 0
	}
	if sec < 60 {
		return fmt.Sprintf("%ds ago", sec)
	}
	min := sec / 60
	if min < 60 {
		return fmt.Sprintf("%dm ago", min)
	}
	hr := min / 60
	if hr < 24 {
		return fmt.Sprintf("%dh ago", hr)
	}
	day := hr / 24
	return fmt.Sprintf("%dd ago", day)
}

// Duration returns a human-friendly duration string between two ISO timestamps.
// If endStr is empty, the current time is used.
func Duration(startStr, endStr string) string {
	if startStr == "" {
		return ""
	}
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return ""
	}
	var end time.Time
	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			return ""
		}
	} else {
		end = time.Now()
	}
	sec := int(math.Floor(end.Sub(start).Seconds()))
	if sec < 0 {
		sec = 0
	}
	if sec < 60 {
		return fmt.Sprintf("%ds", sec)
	}
	min := sec / 60
	remSec := sec % 60
	if min < 60 {
		return fmt.Sprintf("%dm %ds", min, remSec)
	}
	hr := min / 60
	remMin := min % 60
	return fmt.Sprintf("%dh %dm", hr, remMin)
}

// Truncate shortens a string to maxLen, adding an ellipsis if truncated.
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 1 {
		return "\u2026"
	}
	return s[:maxLen-1] + "\u2026"
}

// StatusBadge returns the text and color name for a run/job/step status+conclusion.
func StatusBadge(status types.RunStatus, conclusion types.RunConclusion) (text string, color string) {
	switch status {
	case types.StatusInProgress:
		return "* running", "yellow"
	case types.StatusQueued, types.StatusWaiting, types.StatusPending, types.StatusRequested:
		return "~ queued", "gray"
	}
	// completed — use conclusion
	switch conclusion {
	case types.ConclusionSuccess:
		return "+ passed", "green"
	case types.ConclusionFailure:
		return "x failed", "red"
	case types.ConclusionCancelled:
		return "- cancelled", "gray"
	case types.ConclusionSkipped:
		return "- skipped", "gray"
	case types.ConclusionTimedOut:
		return "! timed out", "red"
	case types.ConclusionActionRequired:
		return "! action req", "yellow"
	default:
		return "? unknown", "gray"
	}
}

// Pad right-pads s with spaces to width.
func Pad(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + spaces(width-len(s))
}

func spaces(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = ' '
	}
	return string(b)
}

package format

import (
	"testing"
	"time"

	"github.com/dzoba/github-actions-watcher/internal/types"
)

func TestRelativeTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		ts   time.Time
		want string
	}{
		{"seconds", now.Add(-30 * time.Second), "30s ago"},
		{"minutes", now.Add(-5 * time.Minute), "5m ago"},
		{"hours", now.Add(-3 * time.Hour), "3h ago"},
		{"days", now.Add(-2 * 24 * time.Hour), "2d ago"},
		{"zero", now, "0s ago"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RelativeTime(tt.ts.Format(time.RFC3339))
			if got != tt.want {
				t.Errorf("RelativeTime() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRelativeTimeInvalid(t *testing.T) {
	if got := RelativeTime("not-a-date"); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestDuration(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	tests := []struct {
		name  string
		start string
		end   string
		want  string
	}{
		{"seconds", base, "2024-01-01T00:00:45Z", "45s"},
		{"minutes", base, "2024-01-01T00:05:30Z", "5m 30s"},
		{"hours", base, "2024-01-01T01:45:00Z", "1h 45m"},
		{"empty start", "", "2024-01-01T00:00:45Z", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Duration(tt.start, tt.end)
			if got != tt.want {
				t.Errorf("Duration() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		s      string
		maxLen int
		want   string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "hell\u2026"},
		{"ab", 1, "\u2026"},
		{"", 5, ""},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := Truncate(tt.s, tt.maxLen)
			if got != tt.want {
				t.Errorf("Truncate(%q, %d) = %q, want %q", tt.s, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestStatusBadge(t *testing.T) {
	tests := []struct {
		status     types.RunStatus
		conclusion types.RunConclusion
		wantText   string
		wantColor  string
	}{
		{types.StatusInProgress, "", "* running", "yellow"},
		{types.StatusQueued, "", "~ queued", "gray"},
		{types.StatusCompleted, types.ConclusionSuccess, "+ passed", "green"},
		{types.StatusCompleted, types.ConclusionFailure, "x failed", "red"},
		{types.StatusCompleted, types.ConclusionCancelled, "- cancelled", "gray"},
		{types.StatusCompleted, types.ConclusionSkipped, "- skipped", "gray"},
		{types.StatusCompleted, types.ConclusionTimedOut, "! timed out", "red"},
		{types.StatusCompleted, types.ConclusionActionRequired, "! action req", "yellow"},
		{types.StatusCompleted, "", "? unknown", "gray"},
	}
	for _, tt := range tests {
		t.Run(tt.wantText, func(t *testing.T) {
			text, color := StatusBadge(tt.status, tt.conclusion)
			if text != tt.wantText || color != tt.wantColor {
				t.Errorf("StatusBadge(%q, %q) = (%q, %q), want (%q, %q)",
					tt.status, tt.conclusion, text, color, tt.wantText, tt.wantColor)
			}
		})
	}
}

func TestPad(t *testing.T) {
	if got := Pad("hi", 5); got != "hi   " {
		t.Errorf("Pad = %q", got)
	}
	if got := Pad("hello", 3); got != "hello" {
		t.Errorf("Pad should not truncate, got %q", got)
	}
}

package ui

import "github.com/charmbracelet/lipgloss"

var (
	CyanBold    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))  // cyan
	Bold        = lipgloss.NewStyle().Bold(true)
	Dim         = lipgloss.NewStyle().Faint(true)
	Blue        = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	Magenta     = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	Yellow      = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	Green       = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	Red         = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	Gray        = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	TabActive   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6")).Background(lipgloss.Color("0"))
	TabInactive = lipgloss.NewStyle().Faint(true)
)

// BadgeStyle returns a lipgloss style for the given color name.
func BadgeStyle(color string) lipgloss.Style {
	switch color {
	case "yellow":
		return Yellow
	case "green":
		return Green
	case "red":
		return Red
	case "gray":
		return Gray
	default:
		return Gray
	}
}

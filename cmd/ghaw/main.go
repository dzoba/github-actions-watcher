package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dzoba/github-actions-watcher/internal/model"
)

func main() {
	interval := flag.Int("i", 10, "Polling interval in seconds")
	flag.IntVar(interval, "interval", 10, "Polling interval in seconds")
	flag.Parse()

	m := model.New(time.Duration(*interval) * time.Second)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

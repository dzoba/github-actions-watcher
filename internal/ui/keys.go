package ui

import "github.com/charmbracelet/bubbles/key"

type ListKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Enter   key.Binding
	Switch  key.Binding
	Refresh key.Binding
	Quit    key.Binding
}

var ListKeys = ListKeyMap{
	Up:      key.NewBinding(key.WithKeys("up", "k")),
	Down:    key.NewBinding(key.WithKeys("down", "j")),
	Enter:   key.NewBinding(key.WithKeys("enter")),
	Switch:  key.NewBinding(key.WithKeys("s")),
	Refresh: key.NewBinding(key.WithKeys("r")),
	Quit:    key.NewBinding(key.WithKeys("q")),
}

type DetailKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Back    key.Binding
	Open    key.Binding
	Refresh key.Binding
	Quit    key.Binding
}

var DetailKeys = DetailKeyMap{
	Up:      key.NewBinding(key.WithKeys("up", "k")),
	Down:    key.NewBinding(key.WithKeys("down", "j")),
	Back:    key.NewBinding(key.WithKeys("esc")),
	Open:    key.NewBinding(key.WithKeys("o")),
	Refresh: key.NewBinding(key.WithKeys("r")),
	Quit:    key.NewBinding(key.WithKeys("q")),
}

type RepoInputKeyMap struct {
	Confirm key.Binding
	Cancel  key.Binding
}

var RepoInputKeys = RepoInputKeyMap{
	Confirm: key.NewBinding(key.WithKeys("enter")),
	Cancel:  key.NewBinding(key.WithKeys("esc")),
}

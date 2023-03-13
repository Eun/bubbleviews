package root

import (
	tea "github.com/charmbracelet/bubbletea"
)

type dummyView struct {
}

func (d *dummyView) Init() tea.Cmd {
	return nil
}

func (d *dummyView) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (d *dummyView) View() string {
	return ""
}

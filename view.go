package bubbleviews

import tea "github.com/charmbracelet/bubbletea"

type ResponseMessage interface {
	View() View
	OnResponse(ResponseMessage) tea.Cmd
}

type View interface {
	Init() tea.Cmd
	Update(tea.Msg) tea.Cmd
	View() string
}

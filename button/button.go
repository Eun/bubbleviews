package button

import (
	"fmt"

	"github.com/Eun/bubbleviews"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ bubbleviews.View = &View{}

type OnResponseFunc func(response *Response) tea.Cmd

type View struct {
	onResponse OnResponseFunc

	style      lipgloss.Style
	focusStyle lipgloss.Style
	text       string
	focusText  string
	focus      bool
	width      int
}

func (m *View) Init() tea.Cmd {
	return nil
}

func (m *View) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m.respond(nil)
		case tea.KeyEsc:
			return m.respond(bubbleviews.EscPressedError{})
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return nil
}

func (m *View) View() string {
	if m.focus {
		return m.focusStyle.MaxWidth(m.width).Render(m.focusText)
	}
	return m.style.MaxWidth(m.width).Render(m.text)
}

func (m *View) Focused() bool {
	return m.focus
}

func (m *View) Focus() {
	m.focus = true
}

func (m *View) Blur() {
	m.focus = false
}

func (m *View) Text() string {
	return m.text
}

func (m *View) SetFocusText(format string, a ...any) {
	m.focusText = fmt.Sprintf(format, a...)
}

func (m *View) FocusText() string {
	return m.focusText
}

func (m *View) SetText(format string, a ...any) {
	m.text = fmt.Sprintf(format, a...)
}

func (m *View) SetStyle(style lipgloss.Style) {
	m.style = style
}

func (m *View) Style() lipgloss.Style {
	return m.style
}

func (m *View) SetFocusStyle(style lipgloss.Style) {
	m.focusStyle = style
}

func (m *View) FocusStyle() lipgloss.Style {
	return m.focusStyle
}

func (m *View) SetOnResponse(fn OnResponseFunc) {
	m.onResponse = fn
}

func (m *View) OnResponse() OnResponseFunc {
	return m.onResponse
}

func (m *View) respond(err error) func() tea.Msg {
	return func() tea.Msg {
		return &Response{
			view:  m,
			Error: err,
		}
	}
}

func New(format string, a ...any) *View {
	var m View
	text := fmt.Sprintf(format, a...)
	m.text = "<" + text + ">"
	m.focusText = "[" + text + "]"
	return &m
}

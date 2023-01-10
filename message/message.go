package message

import (
	"fmt"

	"github.com/Eun/bubbleviews"
	"github.com/Eun/bubbleviews/ext"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
)

var _ bubbleviews.View = &View{}

type View struct {
	OnResponse func(response *Response) tea.Cmd

	viewport viewport.Model
	message  string
	ext.PrefixExt
	ext.SuffixExt
}

func (m *View) Init() tea.Cmd {
	m.viewport.SetYOffset(0)
	return nil
}

func (m *View) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyEsc:
			return m.respond()
		}
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width

		newLineCount := m.PrefixExt.PrefixRenderHeight() + m.SuffixExt.SuffixRenderHeight()
		msg.Height -= newLineCount
		if msg.Height < 0 {
			msg.Height = 0
		}

		m.viewport.Height = msg.Height
		m.viewport.SetContent(wrap.String(wordwrap.String(m.message, msg.Width), msg.Width))
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return cmd
}

func (m *View) View() string {
	return m.RenderPrefix(m.viewport.Width) + m.viewport.View() + m.RenderSuffix(m.viewport.Width)
}

func (m *View) SetMessage(s string) {
	m.message = s
	m.viewport.SetContent(wrap.String(wordwrap.String(m.message, m.viewport.Width), m.viewport.Width))
}

func (m *View) Message() string {
	return m.message
}

func (m *View) respond() func() tea.Msg {
	return func() tea.Msg {
		return &Response{
			model: m,
		}
	}
}

func New(format string, a ...any) *View {
	var m View
	m.message = fmt.Sprintf(format, a...)
	m.viewport = viewport.New(0, 0)
	return &m
}

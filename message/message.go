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

type OnResponseFunc func(response *Response) tea.Cmd

type View struct {
	onResponse OnResponseFunc

	Viewport viewport.Model
	message  string
	ext.PrefixExt
	ext.SuffixExt
}

func (m *View) Init() tea.Cmd {
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
		m.Viewport.Width = msg.Width

		newLineCount := m.PrefixExt.PrefixRenderHeight() + m.SuffixExt.SuffixRenderHeight()
		msg.Height -= newLineCount
		if msg.Height < 0 {
			msg.Height = 0
		}

		m.Viewport.Height = msg.Height
		m.Viewport.SetContent(wrap.String(wordwrap.String(m.message, msg.Width), msg.Width))
	}
	var cmd tea.Cmd
	m.Viewport, cmd = m.Viewport.Update(msg)
	return cmd
}

func (m *View) View() string {
	return m.RenderPrefix(m.Viewport.Width) + m.Viewport.View() + m.RenderSuffix(m.Viewport.Width)
}

func (m *View) SetMessage(s string) {
	m.message = s
	m.Viewport.SetContent(wrap.String(wordwrap.String(m.message, m.Viewport.Width), m.Viewport.Width))
}

func (m *View) Message() string {
	return m.message
}

func (m *View) SetOnResponse(fn OnResponseFunc) {
	m.onResponse = fn
}

func (m *View) OnResponse() OnResponseFunc {
	return m.onResponse
}

func (m *View) respond() func() tea.Msg {
	return func() tea.Msg {
		return &Response{
			view: m,
		}
	}
}

func New(format string, a ...any) *View {
	var m View
	m.message = fmt.Sprintf(format, a...)
	m.Viewport = viewport.New(0, 0)
	return &m
}

package entry

import (
	"github.com/Eun/bubbleviews"
	"github.com/Eun/bubbleviews/ext"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var _ bubbleviews.View = &View{}

type OnResponseFunc func(response *Response) tea.Cmd

type View struct {
	onResponse OnResponseFunc

	ext.PrefixExt
	ext.SuffixExt
	TextInput textinput.Model
}

func (m *View) Init() tea.Cmd {
	return textinput.Blink
}

func (m *View) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.TextInput.Width = msg.Width
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m.respond(nil, bubbleviews.EscPressedError{})
		case tea.KeyEnter:
			v := m.TextInput.Value()
			return m.respond(&v, nil)
		}
	}

	var cmd tea.Cmd
	m.TextInput, cmd = m.TextInput.Update(msg)
	return cmd
}

func (m *View) View() string {
	return m.RenderPrefix(m.TextInput.Width) + m.TextInput.View() + m.RenderSuffix(m.TextInput.Width)
}

func (m *View) SetOnResponse(fn OnResponseFunc) {
	m.onResponse = fn
}

func (m *View) OnResponse() OnResponseFunc {
	return m.onResponse
}

func (m *View) respond(text *string, err error) func() tea.Msg {
	return func() tea.Msg {
		return &Response{
			view:  m,
			Text:  text,
			Error: err,
		}
	}
}

func New() *View {
	var m View
	m.TextInput = textinput.New()
	m.TextInput.Focus()
	return &m
}

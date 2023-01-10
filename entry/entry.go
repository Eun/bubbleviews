package entry

import (
	"github.com/Eun/bubbleviews"
	"github.com/Eun/bubbleviews/ext"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var _ bubbleviews.View = &View{}

type View struct {
	OnResponse func(response *Response) tea.Cmd

	ext.PrefixExt
	ext.SuffixExt
	textinput.Model
}

func (m *View) Init() tea.Cmd {
	m.SetValue("")
	return textinput.Blink
}

func (m *View) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Model.Width = msg.Width
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m.respond(nil, nil)
		case tea.KeyEnter:
			v := m.Model.Value()
			return m.respond(&v, nil)
		}
	}

	var cmd tea.Cmd
	m.Model, cmd = m.Model.Update(msg)
	return cmd
}

func (m *View) View() string {
	return m.RenderPrefix(m.Model.Width) + m.Model.View() + m.RenderSuffix(m.Model.Width)
}

func (m *View) respond(text *string, err error) func() tea.Msg {
	return func() tea.Msg {
		return &Response{
			model: m,
			Text:  text,
			Error: err,
		}
	}
}

func New() *View {
	var m View
	m.Model = textinput.New()
	m.Model.Focus()
	return &m
}

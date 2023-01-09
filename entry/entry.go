package entry

import (
	"strings"

	"github.com/Eun/bubbleviews"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ bubbleviews.View = &View{}

type View struct {
	OnResponse func(response *Response) tea.Cmd

	prefix      string
	prefixStyle lipgloss.Style
	suffix      string
	suffixStyle lipgloss.Style
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
	var sb strings.Builder

	if m.prefix != "" {
		sb.WriteString(m.prefixStyle.MaxWidth(m.Width).Render(m.prefix))
		sb.WriteRune('\n')
	}

	sb.WriteString(m.Model.View())

	if m.suffix != "" {
		sb.WriteRune('\n')
		sb.WriteString(m.suffixStyle.MaxWidth(m.Width).Render(m.suffix))
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (m *View) SetPrefix(s string) {
	m.prefix = s
}

func (m *View) Prefix() string {
	return m.prefix
}

func (m *View) SetPrefixStyle(style lipgloss.Style) {
	m.prefixStyle = style
}

func (m *View) PrefixStyle() lipgloss.Style {
	return m.prefixStyle
}

func (m *View) SetSuffix(s string) {
	m.suffix = s
}

func (m *View) Suffix() string {
	return m.suffix
}

func (m *View) SetSuffixStyle(style lipgloss.Style) {
	m.suffixStyle = style
}

func (m *View) SuffixStyle() lipgloss.Style {
	return m.suffixStyle
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
	m.Model.Width = bubbleviews.Width
	m.Model.Focus()
	return &m
}

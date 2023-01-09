package message

import (
	"fmt"
	"strings"

	"github.com/Eun/bubbleviews"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
)

var _ bubbleviews.View = &View{}

type View struct {
	OnResponse func(response *Response) tea.Cmd

	viewport    viewport.Model
	message     string
	prefix      string
	prefixStyle lipgloss.Style
	suffix      string
	suffixStyle lipgloss.Style
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

		newLineCount := strings.Count(m.viewPrefix()+m.viewSuffix(), "\n")
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

func (m *View) viewPrefix() string {
	if m.prefix == "" {
		return ""
	}
	return m.prefixStyle.MaxWidth(m.viewport.Width).Render(m.prefix) + "\n"
}

func (m *View) viewSuffix() string {
	if m.suffix == "" {
		return ""
	}
	return "\n" + m.suffixStyle.MaxWidth(m.viewport.Width).Render(m.suffix) + "\n"
}

func (m *View) View() string {
	return m.viewPrefix() + m.viewport.View() + m.viewSuffix()
}

func (m *View) SetMessage(s string) {
	m.message = s
	m.viewport.SetContent(wrap.String(wordwrap.String(m.message, m.viewport.Width), m.viewport.Width))
}

func (m *View) Message() string {
	return m.message
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

	m.viewport = viewport.New(bubbleviews.Width, bubbleviews.Height)
	m.viewport.SetContent(wrap.String(wordwrap.String(m.message, bubbleviews.Width), bubbleviews.Width))
	return &m
}

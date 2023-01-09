package loginform

import (
	"strings"

	"github.com/Eun/bubbleviews"
	"github.com/Eun/bubbleviews/button"
	"github.com/Eun/bubbleviews/entry"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ bubbleviews.View = &View{}

type View struct {
	OnResponse func(response *Response) tea.Cmd

	currentFocus bubbleviews.View

	renderer      lipgloss.Style
	prefix        string
	suffix        string
	EntryUsername *entry.View
	EntryPassword *entry.View
	BtnOk         *button.View
	BtnCancel     *button.View

	showOK     bool
	showCancel bool
}

func (m *View) Init() tea.Cmd {
	// set focus
	m.currentFocus = m.EntryUsername
	m.EntryUsername.Focus()
	m.EntryPassword.Blur()
	m.BtnOk.Blur()
	m.BtnCancel.Blur()

	// reset values
	m.EntryUsername.SetValue("")
	m.EntryPassword.SetValue("")

	return tea.Batch(
		m.EntryUsername.Init(),
		m.EntryPassword.Init(),
		m.BtnOk.Init(),
		m.BtnCancel.Init(),
	)
}

func (m *View) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.renderer = m.renderer.MaxWidth(msg.Width)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m.respond(nil, nil, nil)
		case tea.KeyEnter:
			if m.currentFocus == m.BtnOk || (m.currentFocus == m.EntryPassword && !m.showOK) {
				user := m.EntryUsername.Value()
				pass := m.EntryPassword.Value()
				return m.respond(&user, &pass, nil)
			}
			if m.currentFocus == m.BtnCancel {
				return m.respond(nil, nil, nil)
			}
			m.focusNext()
		case tea.KeyTab, tea.KeyDown:
			m.focusNext()
		case tea.KeyShiftTab, tea.KeyUp:
			m.focusPrevious()
		case tea.KeyRight:
			if m.showCancel && m.currentFocus == m.BtnOk {
				m.focusNext()
			}
		case tea.KeyLeft:
			if m.showOK && m.currentFocus == m.BtnCancel {
				m.focusPrevious()
			}
		}
	}

	return tea.Batch(m.EntryUsername.Update(msg), m.EntryPassword.Update(msg))
}

func (m *View) View() string {
	var sb strings.Builder

	if m.prefix != "" {
		sb.WriteString(m.renderer.Render(m.prefix))
		sb.WriteRune('\n')
	}

	sb.WriteString(m.EntryUsername.View())
	sb.WriteRune('\n')
	sb.WriteString(m.EntryPassword.View())
	sb.WriteRune('\n')

	var buttons strings.Builder
	if m.showOK {
		sb.WriteString(m.BtnOk.View())
	}
	if m.showCancel {
		if m.showOK {
			sb.WriteRune(' ')
		}
		sb.WriteString(m.BtnCancel.View())
	}

	if buttons.Len() > 0 {
		sb.WriteString(m.renderer.Render(buttons.String()))
		sb.WriteRune('\n')
	}

	if m.suffix != "" {
		sb.WriteRune('\n')
		sb.WriteString(m.renderer.Render(m.suffix))
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

func (m *View) SetSuffix(s string) {
	m.suffix = s
}

func (m *View) Suffix() string {
	return m.suffix
}

func (m *View) SetShowOK(show bool) {
	m.showOK = show
}

func (m *View) ShowOK() bool {
	return m.showOK
}

func (m *View) SetShowCancel(show bool) {
	m.showCancel = show
}

func (m *View) ShowCancel() bool {
	return m.showCancel
}

func (m *View) respond(username, password *string, err error) func() tea.Msg {
	return func() tea.Msg {
		return &Response{
			model:    m,
			Username: username,
			Password: password,
			Error:    err,
		}
	}
}

func (m *View) focusNext() {
	switch m.currentFocus {
	case m.EntryUsername:
		m.EntryUsername.Blur()
		m.EntryPassword.Focus()
		m.currentFocus = m.EntryPassword
	case m.EntryPassword:
		m.EntryPassword.Blur()
		if m.showOK {
			m.BtnOk.Focus()
			m.currentFocus = m.BtnOk
			break
		}
		if m.showCancel {
			m.BtnCancel.Focus()
			m.currentFocus = m.BtnCancel
			break
		}
		m.EntryUsername.Focus()
		m.currentFocus = m.EntryUsername
	case m.BtnOk:
		m.BtnOk.Blur()
		if m.showCancel {
			m.BtnCancel.Focus()
			m.currentFocus = m.BtnCancel
			break
		}
		m.EntryUsername.Focus()
		m.currentFocus = m.EntryUsername
	case m.BtnCancel:
		m.BtnCancel.Blur()
		m.EntryUsername.Focus()
		m.currentFocus = m.EntryUsername
	}
}

func (m *View) focusPrevious() {
	switch m.currentFocus {
	case m.EntryUsername:
		m.EntryUsername.Blur()
		if m.showCancel {
			m.BtnCancel.Focus()
			m.currentFocus = m.BtnCancel
			break
		}
		if m.showOK {
			m.BtnOk.Focus()
			m.currentFocus = m.BtnOk
			break
		}
		m.EntryPassword.Focus()
		m.currentFocus = m.EntryPassword
	case m.EntryPassword:
		m.EntryPassword.Blur()
		m.EntryUsername.Focus()
		m.currentFocus = m.EntryUsername
	case m.BtnOk:
		m.BtnOk.Blur()
		m.EntryPassword.Focus()
		m.currentFocus = m.EntryPassword
	case m.BtnCancel:
		m.BtnCancel.Blur()
		if m.showOK {
			m.BtnOk.Focus()
			m.currentFocus = m.BtnOk
			break
		}
		m.EntryPassword.Focus()
		m.currentFocus = m.EntryPassword
	}
}

func New() *View {
	var m View
	m.renderer = lipgloss.NewStyle().MaxWidth(bubbleviews.Width)
	m.EntryUsername = entry.New()
	m.EntryUsername.SetPrefix("Username")
	m.EntryPassword = entry.New()
	m.EntryPassword.SetPrefix("Password")
	m.EntryPassword.EchoMode = textinput.EchoPassword

	m.BtnOk = button.New("OK")
	m.BtnCancel = button.New("Cancel")

	return &m
}

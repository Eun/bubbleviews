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
	entryUsername *entry.View
	entryPassword *entry.View
	btnOk         *button.View
	btnCancel     *button.View

	showOK     bool
	showCancel bool
}

func (m *View) Init() tea.Cmd {
	// set focus
	m.currentFocus = m.entryUsername
	m.entryUsername.Focus()
	m.entryPassword.Blur()
	m.btnOk.Blur()
	m.btnCancel.Blur()

	// reset values
	m.entryUsername.SetValue("")
	m.entryPassword.SetValue("")

	return tea.Batch(
		m.entryUsername.Init(),
		m.entryPassword.Init(),
		m.btnOk.Init(),
		m.btnCancel.Init(),
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
			if m.currentFocus == m.btnOk || (m.currentFocus == m.entryPassword && !m.showOK) {
				user := m.entryUsername.Value()
				pass := m.entryPassword.Value()
				return m.respond(&user, &pass, nil)
			}
			if m.currentFocus == m.btnCancel {
				return m.respond(nil, nil, nil)
			}
			m.focusNext()
		case tea.KeyTab, tea.KeyDown:
			m.focusNext()
		case tea.KeyShiftTab, tea.KeyUp:
			m.focusPrevious()
		case tea.KeyRight:
			if m.showCancel && m.currentFocus == m.btnOk {
				m.focusNext()
			}
		case tea.KeyLeft:
			if m.showOK && m.currentFocus == m.btnCancel {
				m.focusPrevious()
			}
		}
	}

	return tea.Batch(m.entryUsername.Update(msg), m.entryPassword.Update(msg))
}

func (m *View) View() string {
	var sb strings.Builder

	if m.prefix != "" {
		sb.WriteString(m.renderer.Render(m.prefix))
		sb.WriteRune('\n')
	}

	sb.WriteString(m.entryUsername.View())
	sb.WriteRune('\n')
	sb.WriteString(m.entryPassword.View())
	sb.WriteRune('\n')

	var buttons strings.Builder
	if m.showOK {
		sb.WriteString(m.btnOk.View())
	}
	if m.showCancel {
		if m.showOK {
			sb.WriteRune(' ')
		}
		sb.WriteString(m.btnCancel.View())
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
	case m.entryUsername:
		m.entryUsername.Blur()
		m.entryPassword.Focus()
		m.currentFocus = m.entryPassword
	case m.entryPassword:
		m.entryPassword.Blur()
		if m.showOK {
			m.btnOk.Focus()
			m.currentFocus = m.btnOk
			break
		}
		if m.showCancel {
			m.btnCancel.Focus()
			m.currentFocus = m.btnCancel
			break
		}
		m.entryUsername.Focus()
		m.currentFocus = m.entryUsername
	case m.btnOk:
		m.btnOk.Blur()
		if m.showCancel {
			m.btnCancel.Focus()
			m.currentFocus = m.btnCancel
			break
		}
		m.entryUsername.Focus()
		m.currentFocus = m.entryUsername
	case m.btnCancel:
		m.btnCancel.Blur()
		m.entryUsername.Focus()
		m.currentFocus = m.entryUsername
	}
}

func (m *View) focusPrevious() {
	switch m.currentFocus {
	case m.entryUsername:
		m.entryUsername.Blur()
		if m.showCancel {
			m.btnCancel.Focus()
			m.currentFocus = m.btnCancel
			break
		}
		if m.showOK {
			m.btnOk.Focus()
			m.currentFocus = m.btnOk
			break
		}
		m.entryPassword.Focus()
		m.currentFocus = m.entryPassword
	case m.entryPassword:
		m.entryPassword.Blur()
		m.entryUsername.Focus()
		m.currentFocus = m.entryUsername
	case m.btnOk:
		m.btnOk.Blur()
		m.entryPassword.Focus()
		m.currentFocus = m.entryPassword
	case m.btnCancel:
		m.btnCancel.Blur()
		if m.showOK {
			m.btnOk.Focus()
			m.currentFocus = m.btnOk
			break
		}
		m.entryPassword.Focus()
		m.currentFocus = m.entryPassword
	}
}

func New() *View {
	var m View
	m.renderer = lipgloss.NewStyle().MaxWidth(bubbleviews.Width)
	m.entryUsername = entry.New()
	m.entryUsername.SetPrefix("Username")
	m.entryPassword = entry.New()
	m.entryPassword.SetPrefix("Password")
	m.entryPassword.EchoMode = textinput.EchoPassword

	m.btnOk = button.New("OK")
	m.btnCancel = button.New("Cancel")

	return &m
}

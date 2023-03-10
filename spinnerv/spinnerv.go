package spinnerv

import (
	"context"

	"github.com/Eun/bubbleviews"
	"github.com/Eun/bubbleviews/ext"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ bubbleviews.View = &View{}

type OnResponseFunc func(response *Response) tea.Cmd

type Action func(ctx context.Context, spinner *View) error

type View struct {
	onResponse OnResponseFunc

	width      int
	widthStyle lipgloss.Style

	Spinner         spinner.Model
	message         string
	messageStyle    lipgloss.Style
	allowEscapeKey  bool
	action          Action
	actionCtx       context.Context
	actionCtxCancel context.CancelFunc
	ext.PrefixExt
	ext.SuffixExt
}

type actionCompleted struct {
	error error
}

func (m *View) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			if m.action == nil {
				return nil
			}
			m.actionCtx, m.actionCtxCancel = context.WithCancel(context.Background())
			return actionCompleted{
				error: m.action(m.actionCtx, m),
			}
		},
		m.Spinner.Tick,
	)
}

func (m *View) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			if m.allowEscapeKey {
				if m.actionCtxCancel != nil {
					m.actionCtxCancel()
				}
				return m.respond(bubbleviews.EscPressedError{})
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.widthStyle = lipgloss.NewStyle().MaxWidth(m.width)
	case actionCompleted:
		select {
		case <-m.actionCtx.Done():
		default:
			return m.respond(msg.error)
		}
	}
	var cmd tea.Cmd
	m.Spinner, cmd = m.Spinner.Update(msg)
	return cmd
}

func (m *View) View() string {
	return m.RenderPrefix(m.width) +
		m.widthStyle.Render(m.Spinner.View()+m.messageStyle.Render(m.message)) +
		m.RenderSuffix(m.width)
}

func (m *View) SetMessage(s string) {
	m.message = s
}

func (m *View) Message() string {
	return m.message
}

func (m *View) SetMessageStyle(s lipgloss.Style) {
	m.messageStyle = s
}

func (m *View) MessageStyle() lipgloss.Style {
	return m.messageStyle
}

func (m *View) SetSpinnerStyle(s lipgloss.Style) {
	m.Spinner.Style = s
}

func (m *View) SpinnerStyle() lipgloss.Style {
	return m.Spinner.Style
}

func (m *View) SetSpinnerType(s spinner.Spinner) {
	m.Spinner.Spinner = s
}

func (m *View) SpinnerType() spinner.Spinner {
	return m.Spinner.Spinner
}

func (m *View) SetAllowEscapeKey(b bool) {
	m.allowEscapeKey = b
}

func (m *View) AllowEscapeKey() bool {
	return m.allowEscapeKey
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
			model: m,
			Error: err,
		}
	}
}

func New(message string, action Action) *View {
	var m View
	m.message = message
	m.Spinner = spinner.New()
	m.action = action
	return &m
}

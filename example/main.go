package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/Eun/bubbleviews"
	"github.com/Eun/bubbleviews/button"
	"github.com/Eun/bubbleviews/entry"
	"github.com/Eun/bubbleviews/example/views/loginform"
	"github.com/Eun/bubbleviews/example/views/selectview"
	"github.com/Eun/bubbleviews/message"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = &TUI{}

var escStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

const (
	maxWidth  = 60
	maxHeight = 20
)

type TUI struct {
	currentModel bubbleviews.View
	quitting     bool

	selectView *selectview.View
}

func (tui *TUI) Init() tea.Cmd {
	return tui.showView(tui.currentModel)
}

func (tui *TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width > maxWidth {
			msg.Width = maxWidth
		}
		if msg.Height > maxHeight {
			msg.Height = maxHeight
		}

		bubbleviews.Width = msg.Width
		bubbleviews.Height = msg.Height
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			tui.quitting = true
			return tui, tea.Quit
		}
	case bubbleviews.ResponseMessage:
		return tui, msg.OnResponse(msg)
	}

	return tui, tui.currentModel.Update(msg)
}

func (tui *TUI) View() string {
	s := tui.currentModel.View()
	if tui.quitting {
		return s + "\n"
	}
	return s
}

func (tui *TUI) showView(model bubbleviews.View) tea.Cmd {
	tui.currentModel = model
	return tea.Batch(tui.currentModel.Init(), func() tea.Msg {
		return tea.WindowSizeMsg{
			Width:  bubbleviews.Width,
			Height: bubbleviews.Height,
		}
	})
}

func (tui *TUI) handleResponse(response interface{}) tea.Cmd {
	msg := message.New("")

	var sb strings.Builder
	enc := json.NewEncoder(&sb)
	enc.SetIndent("", "\t")

	if err := enc.Encode(response); err != nil {
		sb.Reset()
		msg.SetPrefix("Error")
		msg.SetPrefixStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")))
		msg.SetMessage(err.Error())
	} else {
		msg.SetPrefix("Response")
		msg.SetPrefixStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("170")))
		msg.SetMessage(sb.String())
	}

	msg.SetSuffixStyle(escStyle)
	msg.SetSuffix("(esc to go back)")
	msg.OnResponse = func(response *message.Response) tea.Cmd {
		return tui.showView(tui.selectView)
	}
	return tui.showView(msg)
}

func NewTUI() (*TUI, error) { //nolint: unparam // allow nil error
	var tui TUI

	// message view
	msgView := message.New("Hello World")
	msgView.SetSuffixStyle(escStyle)
	msgView.SetSuffix("(esc to go back)")
	msgView.OnResponse = func(response *message.Response) tea.Cmd {
		return tui.handleResponse(response)
	}

	// button view
	buttonView := button.New("Hello World")
	buttonView.OnResponse = func(response *button.Response) tea.Cmd {
		return tui.handleResponse(response)
	}

	// entry view
	entryView := entry.New()
	entryView.SetSuffixStyle(escStyle)
	entryView.SetSuffix("(esc to go back)")
	entryView.OnResponse = func(response *entry.Response) tea.Cmd {
		return tui.handleResponse(response)
	}

	// login view
	loginView := loginform.New()
	loginView.SetShowOK(true)
	loginView.SetShowCancel(true)
	loginView.OnResponse = func(response *loginform.Response) tea.Cmd {
		return tui.handleResponse(response)
	}

	tui.selectView = selectview.New(
		msgView,
		buttonView,
		entryView,
		loginView,
	)
	tui.selectView.OnResponse = func(response *selectview.Response) tea.Cmd {
		if response.SelectedView == nil {
			return nil
		}
		return tui.showView(response.SelectedView)
	}
	tui.currentModel = tui.selectView

	return &tui, nil
}

func main() {
	tui, err := NewTUI()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tea.NewProgram(tui).Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/Eun/bubbleviews"
	"github.com/Eun/bubbleviews/button"
	"github.com/Eun/bubbleviews/entry"
	"github.com/Eun/bubbleviews/example/views/selectview"
	"github.com/Eun/bubbleviews/listv"
	"github.com/Eun/bubbleviews/loginform"
	"github.com/Eun/bubbleviews/message"
	"github.com/Eun/bubbleviews/spinnerv"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sanity-io/litter"
)

var _ tea.Model = &TUI{}

var escStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

const (
	maxWidth  = 60
	maxHeight = 20
)

//go:embed main.go
var f embed.FS

type TUI struct {
	currentModel bubbleviews.View
	quitting     bool

	width  int
	height int
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
		tui.width = msg.Width
		tui.height = msg.Height

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
	return tea.Batch(tea.ClearScreen, tui.currentModel.Init(), func() tea.Msg {
		return tea.WindowSizeMsg{
			Width:  tui.width,
			Height: tui.height,
		}
	})
}

func (tui *TUI) handleResponse(response interface{}) tea.Cmd {
	msg := message.New("")

	msg.SetPrefix("Response")
	msg.SetPrefixStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("170")))
	msg.SetMessage(litter.Sdump(response))

	msg.SetSuffixStyle(escStyle)
	msg.SetSuffix("(esc to go back)")
	msg.SetOnResponse(func(response *message.Response) tea.Cmd {
		return tui.showView(tui.rootView())
	})
	return tui.showView(msg)
}

func (tui *TUI) rootView() bubbleviews.View {
	// message view
	msgView := message.New("")
	data, _ := f.ReadFile("main.go")
	msgView.SetMessage(string(data))

	msgView.SetPrefix("main.go")
	msgView.SetPrefixStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("170")))
	msgView.SetSuffix("(esc to go back)")
	msgView.SetSuffixStyle(escStyle)
	msgView.SetOnResponse(func(response *message.Response) tea.Cmd {
		return tui.handleResponse(response)
	})

	// button view
	buttonView := button.New("Hello World")
	buttonView.SetOnResponse(func(response *button.Response) tea.Cmd {
		return tui.handleResponse(response)
	})

	// entry view
	entryView := entry.New()
	entryView.SetPrefix("Enter some Text")
	entryView.SetPrefixStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("170")))
	entryView.SetSuffix("(esc to go back)")
	entryView.SetSuffixStyle(escStyle)
	entryView.SetOnResponse(func(response *entry.Response) tea.Cmd {
		return tui.handleResponse(response)
	})

	// login view
	loginFormView := loginform.New()
	loginFormView.SetShowOK(true)
	loginFormView.SetShowCancel(true)
	loginFormView.BtnOk.SetFocusStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("170")))
	loginFormView.BtnCancel.SetFocusStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("170")))
	loginFormView.SetPrefix("Please Login")
	loginFormView.SetPrefixStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("170")))
	loginFormView.SetSuffix("(esc to go back)")
	loginFormView.SetSuffixStyle(escStyle)
	loginFormView.SetOnResponse(func(response *loginform.Response) tea.Cmd {
		return tui.handleResponse(response)
	})

	// spinner view
	spinnerView := spinnerv.New(" Loading...", func(ctx context.Context, spinner *spinnerv.View) error {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
		for i := 5; i >= 0; i-- {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(time.Second):
				spinner.SetMessage(fmt.Sprintf(" Loading...%s", style.Render(fmt.Sprintf("%d", i))))
			}
		}
		return nil
	})
	spinnerView.SetPrefix("Please be Patient")
	spinnerView.SetPrefixStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("170")))
	spinnerView.SetSuffix("(esc to go back)")
	spinnerView.SetSuffixStyle(escStyle)
	spinnerView.SetSpinnerStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("170")))
	spinnerView.SetAllowEscapeKey(true)
	spinnerView.SetOnResponse(func(response *spinnerv.Response) tea.Cmd {
		return tui.handleResponse(response)
	})

	items := []listv.ListItem{
		listv.NewSimpleListItem("opt1", "Option 1", ""),
		listv.NewSimpleListItem("opt2", "Option 2", ""),
		listv.NewSimpleListItem("opr3", "Option 3", ""),
	}
	// list view
	listView := listv.New(items,
		listv.NewSimpleListItemDelegate(
			lipgloss.NewStyle().PaddingLeft(2),
			"> ",
			lipgloss.NewStyle().Foreground(lipgloss.Color("170")),
		),
	)
	listView.SetPrefix("Please be Patient")
	listView.SetPrefixStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("170")))
	listView.SetSuffix("(esc to go back)")
	listView.SetSuffixStyle(escStyle)
	listView.SetOnResponse(func(response *listv.Response) tea.Cmd {
		return tui.handleResponse(response)
	})

	rootView := selectview.New(
		msgView,
		buttonView,
		entryView,
		loginFormView,
		spinnerView,
		listView,
	)
	rootView.SetOnResponse(func(response *selectview.Response) tea.Cmd {
		if response.SelectedView == nil {
			return nil
		}
		return tui.showView(response.SelectedView)
	})
	return rootView
}

func NewTUI() (*TUI, error) { //nolint: unparam // allow nil error
	var tui TUI
	tui.showView(tui.rootView())
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

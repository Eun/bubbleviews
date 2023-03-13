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
	"github.com/Eun/bubbleviews/root"
	"github.com/Eun/bubbleviews/spinnerv"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sanity-io/litter"
)

var escStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})
var titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))

//go:embed main.go
var f embed.FS

type RootView struct {
	*root.Root
}

func (rv *RootView) HandleResponse(response interface{}) tea.Cmd {
	msg := message.New("")

	msg.SetPrefix("Response")
	msg.SetPrefixStyle(titleStyle)
	msg.SetMessage(litter.Sdump(response))

	msg.SetSuffixStyle(escStyle)
	msg.SetSuffix("(esc to go back)")
	msg.SetOnResponse(func(response *message.Response) tea.Cmd {
		return rv.SetCurrentView(rv.View())
	})
	return rv.SetCurrentView(msg)
}

func (rv *RootView) View() bubbleviews.View {
	// message view
	msgView := message.New("")
	data, _ := f.ReadFile("main.go")
	msgView.SetMessage(string(data))

	msgView.SetPrefix("main.go")
	msgView.SetPrefixStyle(titleStyle)
	msgView.SetSuffix("(esc to go back)")
	msgView.SetSuffixStyle(escStyle)
	msgView.SetOnResponse(func(response *message.Response) tea.Cmd {
		return rv.HandleResponse(response)
	})

	// button view
	buttonView := button.New("Hello World")
	buttonView.SetOnResponse(func(response *button.Response) tea.Cmd {
		return rv.HandleResponse(response)
	})

	// entry view
	entryView := entry.New()
	entryView.SetPrefix("Enter some Text")
	entryView.SetPrefixStyle(titleStyle)
	entryView.SetSuffix("(esc to go back)")
	entryView.SetSuffixStyle(escStyle)
	entryView.SetOnResponse(func(response *entry.Response) tea.Cmd {
		return rv.HandleResponse(response)
	})

	// login view
	loginFormView := loginform.New()
	loginFormView.SetShowOK(true)
	loginFormView.SetShowCancel(true)
	loginFormView.BtnOk.SetFocusStyle(titleStyle)
	loginFormView.BtnCancel.SetFocusStyle(titleStyle)
	loginFormView.SetPrefix("Please Login")
	loginFormView.SetPrefixStyle(titleStyle)
	loginFormView.SetSuffix("(esc to go back)")
	loginFormView.SetSuffixStyle(escStyle)
	loginFormView.SetOnResponse(func(response *loginform.Response) tea.Cmd {
		return rv.HandleResponse(response)
	})

	// spinner view
	spinnerView := spinnerv.New(" Loading...", func(ctx context.Context, spinner *spinnerv.View) error {
		style := titleStyle
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
	spinnerView.SetPrefix("Please be patient")
	spinnerView.SetPrefixStyle(titleStyle)
	spinnerView.SetSuffix("(esc to go back)")
	spinnerView.SetSuffixStyle(escStyle)
	spinnerView.SetSpinnerStyle(titleStyle)
	spinnerView.SetAllowEscapeKey(true)
	spinnerView.SetOnResponse(func(response *spinnerv.Response) tea.Cmd {
		return rv.HandleResponse(response)
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
			titleStyle,
		),
	)
	listView.SetPrefix("Choose something:")
	listView.SetPrefixStyle(titleStyle)
	listView.SetSuffix("(esc to go back)")
	listView.SetSuffixStyle(escStyle)
	listView.SetOnResponse(func(response *listv.Response) tea.Cmd {
		return rv.HandleResponse(response)
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
		return rv.SetCurrentView(response.SelectedView)
	})
	return rootView
}

func main() {
	rv := RootView{
		Root: root.New(),
	}
	rv.SetCurrentView(rv.View())
	if _, err := tea.NewProgram(rv.Root).Run(); err != nil {
		log.Fatal(err)
	}
}

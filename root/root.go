package root

import (
	"github.com/Eun/bubbleviews"
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &Root{}

type Root struct {
	maxWidth      int
	maxHeight     int
	currentWidth  int
	currentHeight int
	wantView      bubbleviews.View
	currentView   bubbleviews.View
	quitting      bool
}

func (r *Root) MaxWidth() int {
	return r.maxWidth
}

func (r *Root) SetMaxWidth(maxWidth int) {
	r.maxWidth = maxWidth
}

func (r *Root) MaxHeight() int {
	return r.maxHeight
}

func (r *Root) SetMaxHeight(maxHeight int) {
	r.maxHeight = maxHeight
}

func (r *Root) CurrentWidth() int {
	return r.currentWidth
}

func (r *Root) SetCurrentWidth(currentWidth int) {
	r.currentWidth = currentWidth
}

func (r *Root) CurrentHeight() int {
	return r.currentHeight
}

func (r *Root) SetCurrentHeight(currentHeight int) {
	r.currentHeight = currentHeight
}

func (r *Root) CurrentView() bubbleviews.View {
	return r.currentView
}

func (r *Root) SetCurrentView(view bubbleviews.View) tea.Cmd {
	r.wantView = view
	return func() tea.Msg {
		return tea.WindowSizeMsg{
			Width:  r.currentWidth,
			Height: r.currentHeight,
		}
	}
}

func (r *Root) Init() tea.Cmd {
	return nil
}

func (r *Root) View() string {
	s := r.currentView.View()
	if r.quitting {
		return s + "\n"
	}
	return s
}
func (r *Root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width > r.maxWidth {
			msg.Width = r.maxWidth
		}
		if msg.Height > r.maxHeight {
			msg.Height = r.maxHeight
		}
		r.currentWidth = msg.Width
		r.currentHeight = msg.Height

		if r.currentView != r.wantView && r.wantView != nil {
			r.currentView = r.wantView
			return r, tea.Sequence(
				r.currentView.Init(),
				r.currentView.Update(tea.WindowSizeMsg{
					Width:  r.currentWidth,
					Height: r.currentHeight,
				}),
			)
		}

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			r.quitting = true
			return r, tea.Quit
		}
	case bubbleviews.ResponseMessage:
		return r, msg.OnResponse(msg)
	}

	return r, r.currentView.Update(msg)
}

func New() *Root {
	return &Root{
		maxWidth:      60,
		maxHeight:     20,
		currentWidth:  0,
		currentHeight: 0,
		currentView:   &dummyView{},
	}
}

package selectview

import (
	"github.com/Eun/bubbleviews"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ bubbleviews.View = &View{}

type View struct {
	OnResponse func(response *Response) tea.Cmd

	list list.Model
}

func (m *View) Init() tea.Cmd {
	return nil
}

func (m *View) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			i, ok := m.list.SelectedItem().(listItem)
			if ok {
				return m.respond(i.View)
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return cmd
}

func (m *View) View() string {
	return m.list.View()
}

func (m *View) respond(selectedView bubbleviews.View) func() tea.Msg {
	return func() tea.Msg {
		return &Response{
			model:        m,
			SelectedView: selectedView,
		}
	}
}

func New(views ...bubbleviews.View) *View {
	var m View

	items := make([]list.Item, len(views))
	for i := range items {
		items[i] = listItem{
			View: views[i],
		}
	}

	m.list = list.New(items, listItemDelegate{}, bubbleviews.Width, bubbleviews.Height)
	m.list.Title = "Select a View"
	m.list.Styles.Title = lipgloss.NewStyle()
	m.list.Styles.TitleBar = lipgloss.NewStyle()
	m.list.SetShowStatusBar(false)
	m.list.SetShowHelp(false)
	m.list.SetFilteringEnabled(false)
	return &m
}

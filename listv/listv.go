package listv

import (
	"github.com/Eun/bubbleviews"
	"github.com/Eun/bubbleviews/ext"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var _ bubbleviews.View = &View{}

type OnResponseFunc func(response *Response) tea.Cmd
type View struct {
	onResponse OnResponseFunc

	List list.Model
	ext.PrefixExt
	ext.SuffixExt
}

func (m *View) Init() tea.Cmd {
	return nil
}

func (m *View) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m.respond(bubbleviews.EscPressedError{}, nil)
		case tea.KeyEnter:
			item := m.List.SelectedItem()
			if o, ok := item.(ListItem); ok {
				return m.respond(nil, o)
			}
		}
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		m.List.SetHeight(msg.Height - m.PrefixExt.PrefixRenderHeight() - m.SuffixExt.SuffixRenderHeight())
	}
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return cmd
}

func (m *View) View() string {
	return m.RenderPrefix(m.List.Width()) +
		m.List.View() +
		m.RenderSuffix(m.List.Width())
}

func (m *View) SetOnResponse(fn OnResponseFunc) {
	m.onResponse = fn
}

func (m *View) OnResponse() OnResponseFunc {
	return m.onResponse
}

func (m *View) SetItems(items []ListItem) {
	s := make([]list.Item, len(items))
	for i := range items {
		s[i] = items[i]
	}
	m.List.SetItems(s)
}

func (m *View) respond(err error, selection ListItem) func() tea.Msg {
	return func() tea.Msg {
		return &Response{
			view:      m,
			Selection: selection,
			Error:     err,
		}
	}
}

func New(items []ListItem, delegate list.ItemDelegate) *View {
	var m View
	m.List = list.New(nil, delegate, 0, 0)
	m.SetItems(items)
	return &m
}

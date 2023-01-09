package selectview

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	//nolint: gomnd // add 2 spaces
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
)

type listItemDelegate struct{}

func (listItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) { //nolint: gocritic // disable hugeParam
	i, ok := item.(listItem)
	if !ok {
		return
	}

	if index == m.Index() {
		fmt.Fprintf(w, selectedItemStyle.MaxWidth(m.Width()).Render("> %T"), i.View)
		return
	}
	fmt.Fprintf(w, itemStyle.MaxWidth(m.Width()).Render("%T"), i.View)
}

func (listItemDelegate) Height() int {
	return 1
}

func (listItemDelegate) Spacing() int {
	return 0
}

func (listItemDelegate) Update(tea.Msg, *list.Model) tea.Cmd {
	return nil
}

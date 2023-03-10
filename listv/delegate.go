package listv

import (
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type listItemDelegate struct {
	itemStyle          lipgloss.Style
	selectedItemStyle  lipgloss.Style
	selectedItemPrefix string
}

func (d *listItemDelegate) Height() int                               { return 1 }
func (d *listItemDelegate) Spacing() int                              { return 0 }
func (d *listItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d *listItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(ListItem)
	if !ok {
		return
	}

	if index == m.Index() {
		_, _ = io.WriteString(w, d.selectedItemStyle.MaxWidth(m.Width()).Render(d.selectedItemPrefix+i.Title()))
		return
	}
	_, _ = io.WriteString(w, d.itemStyle.MaxWidth(m.Width()).Render(i.Title()))
}

func NewSimpleListItemDelegate(itemStyle lipgloss.Style, selectedItemPrefix string, selectedItemStyle lipgloss.Style) list.ItemDelegate {
	return &listItemDelegate{
		itemStyle:          itemStyle,
		selectedItemStyle:  selectedItemStyle,
		selectedItemPrefix: selectedItemPrefix,
	}
}

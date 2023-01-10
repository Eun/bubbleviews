package ext

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type PrefixExt struct {
	prefixStyle lipgloss.Style
	prefixText  string
}

func (m *PrefixExt) SetPrefixStyle(style lipgloss.Style) {
	m.prefixStyle = style
}

func (m *PrefixExt) PrefixStyle() lipgloss.Style {
	return m.prefixStyle
}

func (m *PrefixExt) SetPrefix(text string) {
	m.prefixText = text
}

func (m *PrefixExt) Prefix() string {
	return m.prefixText
}

func (m *PrefixExt) RenderPrefix(width int) string {
	if m.prefixText == "" {
		return ""
	}
	return m.prefixStyle.MaxWidth(width).Render(m.prefixText) + "\n"
}

func (m *PrefixExt) PrefixRenderHeight() int {
	return strings.Count(m.RenderPrefix(0), "\n")
}

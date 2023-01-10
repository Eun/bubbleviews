package ext

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type SuffixExt struct {
	suffixStyle lipgloss.Style
	suffixText  string
}

func (m *SuffixExt) SetSuffixStyle(style lipgloss.Style) {
	m.suffixStyle = style
}

func (m *SuffixExt) SuffixStyle() lipgloss.Style {
	return m.suffixStyle
}

func (m *SuffixExt) SetSuffix(text string) {
	m.suffixText = text
}

func (m *SuffixExt) Suffix() string {
	return m.suffixText
}

func (m *SuffixExt) RenderSuffix(width int) string {
	if m.suffixText == "" {
		return ""
	}
	return "\n" + m.suffixStyle.MaxWidth(width).Render(m.suffixText) + "\n"
}

func (m *SuffixExt) SuffixRenderHeight() int {
	return strings.Count(m.RenderSuffix(0), "\n")
}

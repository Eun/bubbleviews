package selectview

import (
	"github.com/Eun/bubbleviews"
	tea "github.com/charmbracelet/bubbletea"
)

type Response struct {
	view         *View
	SelectedView bubbleviews.View
}

func (r *Response) View() bubbleviews.View {
	return r.view
}

func (r *Response) OnResponse(msg bubbleviews.ResponseMessage) tea.Cmd {
	response, ok := msg.(*Response)
	if !ok || r.view.onResponse == nil {
		return nil
	}
	return r.view.onResponse(response)
}

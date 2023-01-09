package selectview

import (
	"github.com/Eun/bubbleviews"
	tea "github.com/charmbracelet/bubbletea"
)

type Response struct {
	model        *View
	SelectedView bubbleviews.View
}

func (r *Response) View() bubbleviews.View {
	return r.model
}

func (r *Response) OnResponse(msg bubbleviews.ResponseMessage) tea.Cmd {
	response, ok := msg.(*Response)
	if !ok || r.model.OnResponse == nil {
		return nil
	}
	return r.model.OnResponse(response)
}

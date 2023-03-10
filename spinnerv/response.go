package spinnerv

import (
	"github.com/Eun/bubbleviews"
	tea "github.com/charmbracelet/bubbletea"
)

var _ bubbleviews.ResponseMessage = &Response{}

type Response struct {
	model *View
	Error error
}

func (r *Response) View() bubbleviews.View {
	return r.model
}

func (r *Response) OnResponse(msg bubbleviews.ResponseMessage) tea.Cmd {
	response, ok := msg.(*Response)
	if !ok || r.model.onResponse == nil {
		return nil
	}
	return r.model.onResponse(response)
}

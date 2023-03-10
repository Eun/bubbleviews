package message

import (
	"github.com/Eun/bubbleviews"
	tea "github.com/charmbracelet/bubbletea"
)

var _ bubbleviews.ResponseMessage = &Response{}

type Response struct {
	view *View
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

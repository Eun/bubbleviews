package selectview

import "github.com/Eun/bubbleviews"

type listItem struct {
	bubbleviews.View
}

func (v listItem) FilterValue() string {
	return ""
}

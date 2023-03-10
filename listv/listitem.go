package listv

type ListItem interface {
	ID() string
	Title() string
	FilterValue() string
}

var _ ListItem = &simpleListItem{}

type simpleListItem struct {
	id          string
	title       string
	filterValue string
}

func (s *simpleListItem) ID() string {
	return s.id
}

func (s *simpleListItem) Title() string {
	return s.title
}

func (s *simpleListItem) FilterValue() string {
	return s.filterValue
}

func NewSimpleListItem(id, title, filterValue string) ListItem {
	return &simpleListItem{id: id, title: title, filterValue: filterValue}
}

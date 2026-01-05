package model

type Node struct {
	ID         string  `json:"id"`
	Label      string  `json:"label"`
	NoTick     bool    `json:"noTick"`
	Expandable bool    `json:"expandable"`
	Selectable bool    `json:"selectable"`
	Children   []*Node `json:"children,omitempty"`
	Data       any     `json:"data,omitempty"`
}

func NewNode(id, label string, a ...any) *Node {
	var data any
	if len(a) > 0 {
		data = a[0]
	}
	return &Node{
		ID:         id,
		Label:      label,
		NoTick:     true,
		Expandable: true,
		Selectable: true,
		Data:       data,
	}
}

package model

type Page struct {
	Items []any `json:"items"`
	Size  int   `json:"size"`
	Total int   `json:"total"`
}

func (p Page) IsEmpty() bool {
	return p.Total == 0
}

func (p Page) AddItem(a any) Page {
	p.Items = append(p.Items, a)
	p.Size = len(p.Items)
	return p
}

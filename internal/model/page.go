package model

type Page struct {
	Items []any `json:"items"`
	Size  int   `json:"size"`
	Total int   `json:"total"`
}

func (p Page) IsEmpty() bool {
	return p.Total == 0
}

func (p Page) Add(a any) Page {
	p.Items = append(p.Items, a)
	p.Size++
	return p
}

func (p Page) Set(a []any) Page {
	p.Items = a
	p.Size = len(a)
	return p
}

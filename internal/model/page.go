package model

type Page struct {
	Items []any `json:"items"`
	Size  int   `json:"size"`
	Total int   `json:"total"`
}

func MakePage(items []any, size int, total int) Page {
	if len(items) < size {
		size = len(items)
		total = len(items)
	}
	return Page{Items: items[:size-1], Size: size, Total: total}
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

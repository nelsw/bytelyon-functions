package model

type Page struct {
	Items []any `json:"items"`
	Size  int   `json:"size"`
	Total int   `json:"total"`
}

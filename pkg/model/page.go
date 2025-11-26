package model

import "github.com/oklog/ulid/v2"

type PageType string

const (
	LandingPageType          = "landing"
	ProductPageType          = "product"
	SerpType                 = "serp"
	UnknownPageType PageType = "unknown"
)

type Page struct {
	// ID defines a PK and timestamp
	ID ulid.ULID `json:"id"`
	// Type defines the kind of page
	Type PageType `json:"type"`
	// URL is the page address
	URL string `json:"url"`
	// Title is the <title>
	Title string `json:"title"`
	// Screenshot is a full page .png
	Screenshot []byte `json:"img"`
	// Document is the pages raw HTML
	Document string `json:"html"`
}

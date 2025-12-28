package model

import (
	"errors"
	"regexp"
)

type ProwlerType string

const (
	SearchProwlType  ProwlerType = "search"
	SitemapProwlType ProwlerType = "sitemap"
	ArticleProwlType ProwlerType = "article"
)

var (
	InvalidProwlerType = errors.New("invalid prowler type, must be one of: search, sitemap, article")
)

func (t ProwlerType) String() string {
	return string(t)
}

func (t ProwlerType) Validate() error {
	regex := regexp.MustCompile(`^(search|sitemap|article)$`)
	if !regex.MatchString(t.String()) {
		return InvalidProwlerType
	}
	return nil
}

func NewProwlerType(s string) (t ProwlerType, err error) {
	t = ProwlerType(s)
	return t, t.Validate()
}

func MakeProwlerType(s string) ProwlerType {
	t, err := NewProwlerType(s)
	if err != nil {
		return ""
	}
	return t
}

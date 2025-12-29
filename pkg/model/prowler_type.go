package model

import (
	"errors"
	"regexp"
)

type ProwlerType string

const (
	SearchProwlerType  ProwlerType = "search"
	SitemapProwlerType ProwlerType = "sitemap"
	NewsProwlerType    ProwlerType = "news"
)

var (
	prowlerTypeRegex      = regexp.MustCompile(`^(search|sitemap|news)$`)
	InvalidProwlerTypeErr = errors.New("invalid prowler type, must be one of: search, sitemap, news")
)

func (t ProwlerType) String() string {
	return string(t)
}

func (t ProwlerType) Validate() error {
	if prowlerTypeRegex.MatchString(t.String()) {
		return nil
	}
	return InvalidProwlerTypeErr
}

func NewProwlerType(s string) (ProwlerType, error) {
	return ProwlerType(s), ProwlerType(s).Validate()
}

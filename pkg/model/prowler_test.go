package model

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
)

func TestProwler_Prowl_Search(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db-test")
	p := &Prowler{
		UserID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0"),
		ID:     "indoor bike trainer",
		Type:   SearchProwlerType,
		Targets: Targets{
			"*": true,
		},
	}
	p.Prowl()
}

func TestProwler_Prowl_Sitemap(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db-test")
	p := &Prowler{
		UserID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0"),
		ID:     "https://www.ubicquia.com",
		Type:   SitemapProwlerType,
	}
	p.Prowl()
}

func TestProwler_Prowl_News(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db-test")
	p := &Prowler{
		UserID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0"),
		Type:   NewsProwlerType,
		ID:     "silver price today",
	}
	p.Prowl()
}

func TestProwler_FindAll(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db-test")
	p := &Prowler{
		UserID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0"),
		Type:   NewsProwlerType,
	}
	all, err := p.FindAll()
	assert.NoError(t, err)
	assert.NotEmpty(t, all)
	b, _ := json.MarshalIndent(all, "", "\t")
	fmt.Println(string(b))
}

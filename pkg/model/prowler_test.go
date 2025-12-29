package model

import (
	"encoding/json"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
)

func TestProwler_Prowl_Search(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db-test")
	p := &Prowler{
		UserID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0"),
		ID:     "ev fire blanket",
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
		ID:     "https://www.flowhotel.life",
		Type:   SitemapProwlerType,
	}
	p.Prowl()
}

func TestProwler_Prowl_News(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db-test")
	p := &Prowler{
		UserID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0"),
		Type:   NewsProwlerType,
		ID:     "corsair marine 970",
	}
	p.Prowl()
}

func TestProwler_FindAll(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db-test")
	p := &Prowler{
		UserID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0"),
		Type:   SearchProwlerType,
	}
	all, err := p.FindAll(true)
	assert.NoError(t, err)
	assert.NotEmpty(t, all)
	for _, v := range all {
		b, _ := json.Marshal(v)
		t.Log(string(b))
		assert.NotEmpty(t, v.Results)
	}
}

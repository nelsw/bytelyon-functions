package tests

import (
	. "bytelyon-functions/pkg/model"
	"testing"

	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
)

func init() {
	godotenv.Load()
}

func TestProwler_Prowl_Search(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db-test")
	p := &Prowler{
		UserID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0"),
		ID:     "ev fire blankets",
		Type:   SearchProwlerType,
	}
	p.Prowl()
}

func TestProwler_Prowl_Sitemap(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db-test")
	p := &Prowler{
		UserID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0"),
		ID:     "https://publix.com",
		Type:   SitemapProwlerType,
	}
	p.Prowl()
}

func TestProwler_Prowl_News(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db-test")
	p := &Prowler{
		UserID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0"),
		Type:   NewsProwlerType,
		ID:     "audi r8 reviews",
	}
	p.Prowl()
}

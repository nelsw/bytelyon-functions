package model

import (
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
)

func TestProwler_Prowl(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db")
	p := &Prowler{
		UserID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0"),
		ID:     ulid.MustParse("01KDEWCKTPA7CA6MCRDNZVBRSH"),
		Type:   SearchProwlType,
	}
	for {
		p.Prowl()
		time.Sleep(time.Minute * 5)
	}
}

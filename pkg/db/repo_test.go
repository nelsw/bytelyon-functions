package db

import (
	"testing"

	"github.com/oklog/ulid/v2"
)

func TestMagicDelete(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db-test")
	userID := ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0")
	entityID := ulid.MustParse("01A6JAED00863SNMPBQQ27C39T")
	MagicDelete(userID, entityID)
}

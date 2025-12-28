package model

import (
	"bytelyon-functions/pkg/service/s3"
	"encoding/json"
	"testing"
	"time"
)

func TestProwler_Prowl(t *testing.T) {
	//NewProwler(MakeDemoUser().ID, SearchProwlType, "ev fire blankets", Targets{
	//	"li-fire.com":                 true,
	//	"newpig.com":                  true,
	//	"brimstonefireprotection.com": false,
	//}).Prowl(true)

	var p = new(Prowler)
	b, _ := s3.New().Get("user/01K48PC0BK13BWV2CGWFP8QQH0/prowler/search/01KDEWCKTPA7CA6MCRDNZVBRSH/_.json")
	_ = json.Unmarshal(b, p)

	t.Logf("%+v", p)

	for {
		p.Prowl(true)
		time.Sleep(time.Minute * 5)
	}
}

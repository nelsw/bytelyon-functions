package model

import "testing"

func TestProwler_Prowl(t *testing.T) {

	p := NewProwler(NewUlid(), SearchProwlType, "ev fire blankets", Targets{
		"li-fire.com":                 true,
		"newpig.com":                  true,
		"brimstonefireprotection.com": false,
	})

	p.Prowl(true)
}

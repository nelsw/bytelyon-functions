package pw

import (
	"bytelyon-functions/pkg/model"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ClientSearch(t *testing.T) {

	svc, err := New()
	assert.NoError(t, err)

	var html string
	var img []byte

	html, img, err = svc.Search("INEOS Grenadier for sale in Fort Lauderdale")
	assert.NoError(t, err)

	_ = os.WriteFile("tmp.png", img, os.ModePerm)

	model.MakeHTML(html).Formatted()
}

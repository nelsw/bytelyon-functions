package model

import (
	logger "bytelyon-functions/pkg"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Crawl(t *testing.T) {
	t.Setenv("APP_MODE", "test")
	logger.Init()
	user := NewDemoUser()
	sitemap := NewSitemap(user)
	sitemap.URL = "https://www.ford.com"
	b, _ := json.Marshal(&sitemap)
	out, err := sitemap.Create(b)
	assert.NoError(t, err)
	assert.NotNil(t, out)

	fmt.Println(out)
}

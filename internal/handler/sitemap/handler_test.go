package sitemap

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/handler/jwt"
	"bytelyon-functions/internal/model"
	"bytelyon-functions/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Handler(t *testing.T) {

	req := test.
		NewRequest(t).
		Bearer(jwt.CreateString(test.CTX, test.DemoUser())).
		Post(model.Sitemap{URL: "https://ubicquia.com", Depth: 25})

	res, _ := Handler(test.CTX, req)

	var sitemap model.Sitemap
	app.MustUnmarshal([]byte(res.Body), &sitemap)
	assert.Greater(t, len(sitemap.Tracked), 10)
	assert.Greater(t, len(sitemap.Visited), 10)
}

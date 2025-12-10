package model

import (
	logger "bytelyon-functions/pkg"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
)

func TestFindLastArticle(t *testing.T) {

	user := MakeDemoUser()
	article := &Article{News: &News{
		User: &user,
		ID:   ulid.MustParse("01KB0MB9ZD0Z8MM0P5MVFWE3YN"),
	}}
	err := article.FindLast()
	assert.NoError(t, err)
}

func TestFindArticle(t *testing.T) {
	user := MakeDemoUser()
	a := NewArticle(&user, "01KB0MB9ZD0Z8MM0P5MVFWE3YN", "01KB0MBA556ZJR7JNRK8YHQFTV")
	err := a.Find()
	assert.NoError(t, err)
}

func TestFindAllArticles(t *testing.T) {
	user := MakeDemoUser()
	all, err := NewArticle(&user, "01KB0MB9ZD0Z8MM0P5MVFWE3YN").FindAll()
	assert.NoError(t, err)
	assert.NotEmpty(t, all)
}

func TestArticle_Delete(t *testing.T) {
	t.Setenv("APP_MODE", "test")
	logger.Init()
	user := MakeDemoUser()
	a := NewArticle(&user, "01KB0MB9ZD0Z8MM0P5MVFWE3YN", "01KB0MD9RGTKEXB7DQTV61KRHY")
	err := a.Delete()
	assert.NoError(t, err)
}

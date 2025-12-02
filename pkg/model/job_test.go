package model

import (
	"bytelyon-functions/pkg/util/pretty"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	user := MakeDemoUser()
	job := NewJob(&user)

	t.Log(job)
	pretty.Println(job)

	fmt.Println(time.Now().Format(time.RFC3339))
	fmt.Println(time.Now().UTC().Format(time.RFC3339))
}

func TestJob_Save(t *testing.T) {
	user := MakeDemoUser()
	job, err := NewJob(&user).Save([]byte(`
{
  "id":"01KB0MB9ZD0Z8MM0P5MVFWE3YN",
  "type":"news",
  "frequency": {
    "unit": "m",
    "value": 10
  }
}`))
	assert.NoError(t, err)
	assert.Equal(t, job.User.ID, user.ID)
	assert.Equal(t, job.Type, NewsJobType)

	pretty.Println(job)
}

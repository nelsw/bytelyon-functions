package internal

import (
	"fmt"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	InitLogger()
}

func TestClient_Keys(t *testing.T) {
	out, err := S3.Keys("user/01K48PC0BK13BWV2CGWFP8QQH0/job/")
	assert.NoError(t, err)
	assert.NotEmpty(t, out)
	for _, k := range out {
		log.Debug().Msg(k)
	}
}

func TestClient_Get(t *testing.T) {
	out, err := S3.Get("user/01K48PC0BK13BWV2CGWFP8QQH0/_.json")
	assert.NoError(t, err)
	fmt.Println(string(out))
}

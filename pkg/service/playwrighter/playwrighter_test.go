package playwrighter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ClientSearch(t *testing.T) {

	svc, err := New()
	assert.NoError(t, err)

	err = svc.Search("ev fire blanket")
	assert.NoError(t, err)
}

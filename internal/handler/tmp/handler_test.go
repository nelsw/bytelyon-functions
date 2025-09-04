package tmp

import (
	"bytelyon-functions/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	/*
		Arrange
	*/
	req := test.NewRequest(t).Patch()

	/*
		Act
	*/
	res, err := Handler(req)

	/*
		Assert
	*/
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)
}

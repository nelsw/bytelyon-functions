package tmp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	/*
		Arrange
	*/

	/*
		Act
	*/
	err := Handler(Request{"https://google.com/search?q=corsair+marine+970"})

	/*
		Assert
	*/
	assert.NoError(t, err)
}

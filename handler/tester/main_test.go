package main

import (
	"bytelyon-functions/pkg/api"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	/*
		Arrange
	*/
	req := api.NewRequest().Get()

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

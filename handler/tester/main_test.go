package main

import (
	"bytelyon-functions/pkg/api"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Global setup: e.g., connect to a database, set up a mock server
	fmt.Println("Performing global setup...")
	// Run all tests in the package
	code := m.Run()
	// Global teardown: e.g., close database connections, clean up resources
	fmt.Println("Performing global teardown...")
	os.Exit(code)
}

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

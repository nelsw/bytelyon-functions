package main

import (
	"os"
	"testing"
)

func TestBuildZip(t *testing.T) {
	BuildZip("contact")

	if _, err := os.Stat("bootstrap"); os.IsNotExist(err) {
		t.Error("build failed")
	}

	if _, err := os.Stat("main.zip"); os.IsNotExist(err) {
		t.Error("zip failed")
	}
}

func TestCleanup(t *testing.T) {
	Cleanup()
}

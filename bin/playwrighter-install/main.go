package main

import "github.com/playwright-community/playwright-go"

func main() {
	if err := playwright.Install(); err != nil {
		panic(err)
	}
}

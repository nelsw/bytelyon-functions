package main

import (
	"bufio"
	. "bytelyon-functions/cmd/prowl/internal"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

var w = NewWorker()

func init() {
	godotenv.Load()
}

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go w.Start()

	go func() {
		// Initialize a scanner to read input line by line
		scanner := bufio.NewScanner(os.Stdin)

		for {
			// Scan for the next line of input
			if !scanner.Scan() {
				break // Exit the loop if there's an input error or EOF (e.g., Ctrl+D)
			}

			if strings.TrimSpace(scanner.Text()) == "q" {
				stop()
				return
			}

			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
			}
		}
	}()

	fmt.Println(`🦁 Ready!`)

	<-ctx.Done()

	fmt.Println(`🦁 Exiting...`)

	for w.Stop(); !w.Done(); time.Sleep(1 * time.Second) {
	}

	fmt.Println(`🦁 Goodbye!`)

	os.Exit(0)
}

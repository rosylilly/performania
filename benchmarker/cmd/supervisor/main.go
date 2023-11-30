package main

import (
	"context"
	"os"
	"os/signal"
)

func main() {
	_, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
}

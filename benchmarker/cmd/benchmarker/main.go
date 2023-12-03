package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/isucon/isucandar"
	"github.com/rosylilly/performania/benchmarker"
	"github.com/rosylilly/performania/benchmarker/scenario"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	benchmark, err := isucandar.NewBenchmark(
		isucandar.WithPrepareTimeout(60*time.Second),
		isucandar.WithLoadTimeout(60*time.Second),
		isucandar.WithoutPanicRecover(),
	)
	if err != nil {
		panic(err)
	}

	benchmark.AddScenario(scenario.NewScenario())

	result := benchmark.Start(ctx)

	for _, err := range result.Errors.All() {
		log.Printf("[ERROR] %+v", err)
	}

	for tag, score := range result.Score.Breakdown() {
		benchmarker.UserLogger.Printf("[SCORE] %s: %d", tag, score)
	}
	benchmarker.UserLogger.Printf("[SCORE] Total: %d", result.Score.Total())
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()

	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(textHandler)

	fmt.Println("--- Day One ---")
	resultA, err := OneA(ctx, logger)
	if err != nil {
		logger.Error("OneA failed", slog.String("error", err.Error()))
	}

	resultB, err := OneB(ctx, logger)
	if err != nil {
		logger.Error("OneB failed", slog.String("error", err.Error()))
	}

	fmt.Printf("A: %d\n", resultA)
	fmt.Printf("B: %d\n", resultB)
}

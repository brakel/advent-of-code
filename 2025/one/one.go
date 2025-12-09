package main

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

//go:embed data/input_A.txt
var inputA []byte

//go:embed data/input_B.txt
var inputB []byte

const (
	startingPosition = 50
	maxDegrees       = 100
)

type ImplFunc func(currentPos, direction, zeroCount int) (int, int)

// OneA helps the Elves bypass a decoy safe to find the real password for a secret entrance.
// It calculates the password by simulating the rotations of a safe dial (numbered 0-99, starting at 50)
// and counting how many times it lands on the magic number 0.
func OneA(ctx context.Context, logger *slog.Logger) (int, error) {
	impl := func(currentPos, direction, zeroCount int) (int, int) {
		newPos := (currentPos + direction) % maxDegrees
		if newPos == 0 {
			zeroCount++
		}

		return newPos, zeroCount
	}

	return run(inputA, impl)
}

// OneB implements the "newer" 0x434C49434B security protocol for the secret entrance.
// This method requires a more detailed password calculation. The password is the
// total number of times the safe dial points to 0 during any click, including
// intermediate clicks within a rotation, not just when a rotation finishes.
func OneB(ctx context.Context, logger *slog.Logger) (int, error) {
	impl := func(currentPos, direction, zeroCount int) (int, int) {
		newPos := (currentPos + direction) % maxDegrees
		if newPos == 0 {
			zeroCount++
		}

		if newPos < 0 {
			newPos += maxDegrees
		}

		partialRotationMagnitude := direction % maxDegrees
		partialCrossesZero := currentPos+partialRotationMagnitude > maxDegrees || currentPos+partialRotationMagnitude < 0
		fullRotationCount := abs((direction - partialRotationMagnitude) / maxDegrees)

		logger.DebugContext(
			ctx,
			"Calculating full rotations over zero",
			slog.Int("current_position", currentPos),
			slog.Int("direction", direction),
			slog.Int("new_position", newPos),
			slog.Int("partial_rotation_magnitude", partialRotationMagnitude),
			slog.Bool("partial_crosses_zero", partialCrossesZero),
			slog.Int("full_rotation_count", fullRotationCount),
		)

		if partialCrossesZero && currentPos != 0 {
			fullRotationCount++
		}

		zeroCount += fullRotationCount

		return newPos, zeroCount
	}

	return run(inputB, impl)
}

func run(input []byte, impl ImplFunc) (int, error) {
	i := bytes.NewReader(input)
	scanner := bufio.NewScanner(i)

	currentPos := startingPosition
	newPos := 0
	zeroCount := 0

	for scanner.Scan() {
		code := scanner.Text()
		direction, err := parse(code)
		if err != nil {
			return 0, fmt.Errorf("failed to parse code %s: %w", code, err)
		}

		newPos, zeroCount = impl(currentPos, direction, zeroCount)

		currentPos = newPos
	}

	return zeroCount, nil
}

func parse(code string) (int, error) {
	parsedCode := strings.NewReplacer("R", "", "L", "-").Replace(code)

	direction, err := strconv.Atoi(parsedCode)
	if err != nil {
		return 0, err
	}

	return direction, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

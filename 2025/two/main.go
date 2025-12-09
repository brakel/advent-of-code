package main

import (
	"context"
	_ "embed"
	"fmt"
	"iter"
	"math"
	"strconv"
	"strings"
	"sync"
)

//go:embed data/input_A.txt
var inputA string

type IDRange struct {
	Lower Bound
	Upper Bound
}

type Bound struct {
	Length int
	Value  int
}

func main() {
	ctx := context.Background()
	c := make(chan int)

	var wg sync.WaitGroup

	TwoA(ctx, c, &wg)

	go func() {
		wg.Wait()
		close(c)
	}()

	total := 0
	for t := range c {
		total += t
	}

	fmt.Println("--- Day Two ---")
	fmt.Println("A:", total)
}

// TwoA identifies invalid IDs within given ranges.
// An ID is considered invalid if it is formed by repeating a sequence of digits twice (e.g., 55, 6464, 123123).
// Leading zeros are not allowed for IDs.
// The function processes a comma-separated list of ID ranges,
// and for each range, it finds and sums all invalid IDs.
func TwoA(ctx context.Context, c chan int, wg *sync.WaitGroup) {
	for idr := range parseIDs(inputA) {

		wg.Go(func() {
			lowerOdd := false
			upperOdd := false

			if idr.Lower.Length%2 != 0 {
				lowerOdd = true
			}

			if idr.Upper.Length%2 != 0 {
				upperOdd = true
			}

			// Skip if both bounds are odd and of the same length.
			if lowerOdd && upperOdd {
				if idr.Lower.Length == idr.Upper.Length {
					return
				}
			}

			// If one of the bounds is odd length, we need to see if we can adjust to an even length within bounds.
			if lowerOdd {
				nextEvenLength := idr.Lower.Length + 1
				minValue := int(math.Pow10(idr.Lower.Length))

				idr.Lower.Length = nextEvenLength
				idr.Lower.Value = minValue
			}

			if upperOdd {
				prevEvenLength := idr.Upper.Length - 1
				maxValue := int(math.Pow10(prevEvenLength)) - 1

				idr.Upper.Length = prevEvenLength
				idr.Upper.Value = maxValue
			}

			diff := idr.Upper.Value - idr.Lower.Value

			// Split the lower bound value in half to get the first sequence to match against.
			firstSeq := int(float64(idr.Lower.Value) / math.Pow10(idr.Lower.Length/2))

			// The cycle length if the length the diff needs to exceed for the first half to increment.
			// If the diff is greater than or equal to 10^(length/2), the first half can increment creating new invalid IDs.
			// e.g. 11-22 has a diff of 11, which is geater than 10^(2/2)=10, so "firstSeq" can increment from 1 to 2.
			//
			// Also useful for constructing invalid IDs - see following check.
			cycleLength := int(math.Pow10(idr.Lower.Length / 2))

			// Calculate how many full cycles fit in the range
			fullCycles := diff / cycleLength

			total := 0

			// There will be an extra cycle to check if the diff is not an exact multiple of the cycle length.
			// This can be below the lower bound or above the upper bound, so we need to check.
			for i := 0; i <= fullCycles+1; i++ {
				invalidID := (firstSeq+i)*cycleLength + (firstSeq + i)

				if invalidID < idr.Lower.Value {
					continue
				}

				if invalidID > idr.Upper.Value {
					break
				}

				total += invalidID
			}

			c <- total
		})
	}
}

// idRanges parses a comma-separated string of ID ranges, returning a iterator over IDRange structs.
func parseIDs(input string) iter.Seq[IDRange] {
	idRanges := strings.SplitSeq(input, ",")

	return func(yield func(IDRange) bool) {
		for r := range idRanges {
			r = strings.ReplaceAll(r, "\n", "")

			bounds := strings.SplitN(r, "-", 2)

			lowerBoundLength := len(bounds[0])
			upperBoundLength := len(bounds[1])

			lowerBoundValue, err := strconv.Atoi(bounds[0])
			if err != nil {
				return
			}

			upperBoundValue, err := strconv.Atoi(bounds[1])
			if err != nil {
				return
			}

			idr := IDRange{
				Lower: Bound{
					Length: lowerBoundLength,
					Value:  lowerBoundValue,
				},
				Upper: Bound{
					Length: upperBoundLength,
					Value:  upperBoundValue,
				},
			}

			if !yield(idr) {
				return
			}
		}
	}
}

package main

import (
	"fmt"

	"github.com/gverger/advent2021/utils"
	"github.com/gverger/advent2021/utils/maps"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {
	values, err := maps.Strings(lines).ToInts()
	if err != nil {
		return fmt.Errorf("cannot convert data to list of ints: %w", err)
	}

	count, err := incrCount(values, 3)
	if err != nil {
		return fmt.Errorf("cannot count increases: %w", err)
	}

	fmt.Println("Nb of increases: ", count)
	return nil
}

func incrCount(values []int, windowSize int) (int, error) {
	nbIncr := 0

	for i, current := range values[windowSize:] {
		last := values[i]
		if current > last {
			nbIncr++
		}
	}

	return nbIncr, nil
}

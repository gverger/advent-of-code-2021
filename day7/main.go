package main

import (
	"fmt"
	"strings"

	"github.com/gverger/advent2021/utils"
	"github.com/gverger/advent2021/utils/maps"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {
	positions, err := positionsFromInput(lines[0])
	if err != nil {
		return err
	}

	pos1, score1 := bestPos(positions, part1FuelSpent)
	fmt.Println("Result (part 1):", pos1, score1)

	pos2, score2 := bestPos(positions, part2FuelSpent)
	fmt.Println("Result (part 2):", pos2, score2)

	return nil
}

type Position = int

// bestPos returns the fuel spent and the best position
func bestPos(positions []Position, compute fuelComputation) (int, int) {
	sum := 0
	for _, p := range positions {
		sum += p
	}

	best := sum / len(positions)

	fuel := fuelIfMoveTo(positions, best, compute)

	fuelRight := fuelIfMoveTo(positions, best+1, compute)
	for fuelRight < fuel {
		fuel = fuelRight
		best = best + 1
		fuelRight = fuelIfMoveTo(positions, best+1, compute)
	}

	fuelLeft := fuelIfMoveTo(positions, best-1, compute)
	for fuelLeft < fuel {
		fuel = fuelLeft
		best = best - 1
		fuelLeft = fuelIfMoveTo(positions, best-1, compute)
	}

	return fuel, best
}

type fuelComputation func(from, to Position) int

func fuelIfMoveTo(positions []Position, aim Position, compute fuelComputation) int {
	move := 0
	for _, p := range positions {
		move += compute(p, aim)
	}

	return move
}

func part1FuelSpent(from int, to int) int {
	if from > to {
		return from - to
	}
	return to - from
}

func part2FuelSpent(from int, to int) int {
	n := 0
	if from > to {
		n = from - to
	} else {
		n = to - from
	}

	return n * (n + 1) / 2
}

func positionsFromInput(line string) ([]Position, error) {
	pos := strings.Split(line, ",")

	return maps.Strings(pos).ToInts()
}

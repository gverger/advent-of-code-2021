package main

import (
	"fmt"

	m "github.com/gverger/advent2021/day2/moves"
	"github.com/gverger/advent2021/day2/part1"
	"github.com/gverger/advent2021/day2/part2"
	"github.com/gverger/advent2021/utils"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {
	moves, err := m.FromInput(lines)
	if err != nil {
		return err
	}

	arrival1 := part1.NewPos().Apply(moves)
	fmt.Printf("Result part 2: %d (x=%d, y=%d)\n", arrival1.X*arrival1.Y, arrival1.X, arrival1.Y)

	arrival2 := part2.NewPos().Apply(moves)
	fmt.Printf("Result part 2: %d (x=%d, y=%d, aim=%d)\n", arrival2.X*arrival2.Y, arrival2.X, arrival2.Y, arrival2.Aim)

	return nil
}

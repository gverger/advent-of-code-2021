package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	m "github.com/gverger/advent2021/day2/moves"
	"github.com/gverger/advent2021/day2/part1"
	"github.com/gverger/advent2021/day2/part2"
)

var inputFile = flag.String("input", "input.txt", "the input file")

func main() {
	flag.Parse()

	if inputFile == nil {
		fmt.Println("ERROR: empty input name")
		os.Exit(1)
	}

	err := run(*inputFile)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(2)
	}
}

func run(input string) error {
	lines, err := readLines(input)
	if err != nil {
		return fmt.Errorf("cannot read lines: %w", err)
	}

	if len(lines) == 0 {
		return errors.New("no line")
	}

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

func readLines(fileName string) ([]string, error) {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot read %q: %w", fileName, err)
	}

	data := strings.TrimSpace(string(raw))

	return strings.Split(data, "\n"), nil
}

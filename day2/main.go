package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	moves, err := movesFromInput(lines)
	if err != nil {
		return err
	}

	end := Pos{X: 0, Y: 0}.Apply(moves)

	fmt.Printf("Result: %d (x=%d, y=%d)", end.X*end.Y, end.X, end.Y)

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

type Pos struct {
	X int
	Y int
}

func (p Pos) Apply(moves []Move) Pos {
	newPos := Pos{X: p.X, Y: p.Y}

	for _, m := range moves {
		switch m.Dir {
		case Up:
			newPos.Y -= m.Units
		case Down:
			newPos.Y += m.Units
		case Forward:
			newPos.X += m.Units
		}
	}

	return newPos
}

func movesFromInput(lines []string) ([]Move, error) {
	moves := make([]Move, 0, len(lines))

	for _, l := range lines {
		m, err := moveFromInput(l)
		if err != nil {
			return moves, err
		}
		moves = append(moves, m)
	}

	return moves, nil
}

type Direction string

const (
	Forward Direction = "forward"
	Up      Direction = "up"
	Down    Direction = "down"
)

func (dir Direction) IsValid() bool {
	return dir == Forward || dir == Up || dir == Down
}

type Move struct {
	Dir   Direction
	Units int
}

func moveFromInput(line string) (Move, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 2 {
		return Move{}, fmt.Errorf("not a valid input: %q", line)
	}
	dir := Direction(parts[0])
	if !dir.IsValid() {
		return Move{}, fmt.Errorf("not a valid direction: %q", dir)
	}
	units, err := strconv.Atoi(parts[1])
	if err != nil {
		return Move{}, err
	}
	return Move{Dir: dir, Units: units}, nil
}

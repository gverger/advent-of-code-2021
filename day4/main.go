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

	choices, err := choicesFrom(lines[0])
	if err != nil {
		return err
	}

	boards, err := boardsFrom(lines[2:])
	if err != nil {
		return err
	}

	fmt.Println("Best score (part 1):", part1Score(boards, choices))
	fmt.Println("Worst score (part 2):", part2Score(boards, choices))

	return nil
}

func part1Score(boards []Board, choices []int) int {
	firstDone := len(choices) + 1
	bestScore := 0

	for _, b := range boards {
		doneAt, score := b.score(choices)
		if doneAt < firstDone {
			bestScore = score
			firstDone = doneAt
		}
	}

	return bestScore
}

func part2Score(boards []Board, choices []int) int {
	firstDone := -1
	worstScore := 0

	for _, b := range boards {
		doneAt, score := b.score(choices)
		if doneAt > firstDone {
			worstScore = score
			firstDone = doneAt
		}
	}

	return worstScore
}

func readLines(fileName string) ([]string, error) {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot read %q: %w", fileName, err)
	}

	data := strings.TrimSpace(string(raw))

	return strings.Split(data, "\n"), nil
}

func choicesFrom(line string) ([]int, error) {
	inputs := strings.Split(line, ",")

	choices := make([]int, 0)

	for _, input := range inputs {
		choice, err := strconv.Atoi(input)
		if err != nil {
			return nil, fmt.Errorf("wrong draw: %q", input)
		}
		choices = append(choices, choice)
	}

	return choices, nil
}

type Pos struct {
	X int
	Y int
}

type Board struct {
	numbers map[int]Pos
	lines   [][]int
}

func (b Board) String() string {
	var builder strings.Builder
	for _, lines := range b.lines {
		for _, n := range lines {
			builder.WriteString(strconv.Itoa(n))
			builder.WriteString(" ")
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

func (b Board) totalSum() int {
	s := 0
	for n := range b.numbers {
		s += n
	}
	return s
}

func (b *Board) score(choices []int) (int, int) {
	colScores := make([]int, 5)
	rowScores := make([]int, 5)

	sumFound := 0
	lastStep := -1
	lastChoice := -1

	for step, c := range choices {
		pos, found := b.numbers[c]
		if !found {
			continue
		}
		sumFound += c
		colScores[pos.X] += 1
		rowScores[pos.Y] += 1
		if colScores[pos.X] == 5 || rowScores[pos.Y] == 5 {
			lastStep = step
			lastChoice = c
			break
		}
	}

	return lastStep, (b.totalSum() - sumFound) * lastChoice
}

func NewBoard(lines []string) (Board, error) {
	b := Board{
		lines:   make([][]int, 5),
		numbers: make(map[int]Pos),
	}

	for i, l := range lines {
		numbers, err := StringSlice(strings.Fields(l)).mapToInts()
		if err != nil {
			return Board{}, fmt.Errorf("cannot get numbers for %q: %w", l, err)
		}

		for idxN, n := range numbers {
			b.numbers[n] = Pos{X: idxN, Y: i}
		}

		b.lines[i] = numbers
	}

	return b, nil
}

func boardsFrom(lines []string) ([]Board, error) {
	boards := make([]Board, 0)

	i := 0
	for i < len(lines) {
		currentBoard, err := NewBoard(lines[i:(i + 5)])
		if err != nil {
			return nil, err
		}

		boards = append(boards, currentBoard)
		i += 6
	}

	return boards, nil
}

type StringSlice []string

func (data StringSlice) mapToInts() ([]int, error) {
	values := make([]int, 0, len(data))
	for _, text := range data {
		current, err := strconv.Atoi(text)
		if err != nil {
			return nil, err
		}

		values = append(values, current)
	}

	return values, nil
}

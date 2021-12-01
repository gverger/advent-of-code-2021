package main

import (
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

	values, err := StringSlice(lines).mapToInts()
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

func readLines(fileName string) ([]string, error) {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot read %q: %w", fileName, err)
	}

	data := strings.TrimSpace(string(raw))

	return strings.Split(data, "\n"), nil
}

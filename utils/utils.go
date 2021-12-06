package utils

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

var inputFile = flag.String("input", "input.txt", "the input file")

func Main(fn func(lines []string) error) {
	flag.Parse()

	if inputFile == nil {
		fmt.Println("ERROR: empty input name")
		os.Exit(1)
	}

	err := Run(*inputFile, fn)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(2)
	}
}

func ReadLines(fileName string) ([]string, error) {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot read %q: %w", fileName, err)
	}

	data := strings.TrimSpace(string(raw))

	return strings.Split(data, "\n"), nil
}

func Run(fileName string, fn func(lines []string) error) error {
	lines, err := ReadLines(fileName)
	if err != nil {
		return fmt.Errorf("cannot read lines: %w", err)
	}

	if len(lines) == 0 {
		return errors.New("no line")
	}

	return fn(lines)
}

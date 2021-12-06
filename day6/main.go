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

	timers := timersFromInput(lines[0])

	school := NewSchool()
	for _, t := range timers {
		school.AddFishes(t, 1)
	}

	for i := 0; i < 256; i++ {
		school.NextDay()
	}

	fmt.Println(school.NbFishes())

	return nil
}

type School struct {
	timers     []int64
	currentDay int
}

func NewSchool() School {
	return School{
		timers:     make([]int64, 9),
		currentDay: 0,
	}
}

func (s School) NbFishes() int64 {
	total := int64(0)
	for _, nb := range s.timers {
		total += int64(nb)
	}

	return total
}

func (s *School) AddFishes(timer int, nb int64) {
	s.timers[(timer+s.currentDay)%9] += nb
}

func (s *School) NextDay() {
	nbFishes := s.timers[s.currentDay]
	// If the list was not circular of size 9 we should do
	// --
	//    s.timers[s.currentDay] = 0
	//    s.AddFishes(9, nbFishes)
	// --
	s.AddFishes(7, nbFishes)
	s.currentDay += 1
	if s.currentDay >= 9 {
		s.currentDay = 0
	}
}

func timersFromInput(input string) []int {
	numbers := strings.Split(input, ",")

	timers := make([]int, len(numbers))
	for i, s := range numbers {
		n, _ := strconv.Atoi(s)
		timers[i] = n
	}

	return timers
}

func readLines(fileName string) ([]string, error) {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot read %q: %w", fileName, err)
	}

	data := strings.TrimSpace(string(raw))

	return strings.Split(data, "\n"), nil
}

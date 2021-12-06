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
	timers, _ := maps.Strings(numbers).ToInts()

	return timers
}

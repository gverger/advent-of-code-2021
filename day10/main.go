package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/gverger/advent2021/utils"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {

	fmt.Println("Syntax Error cost:", part1(lines))
	fmt.Println("Completion cost:", part2(lines))

	return nil
}

func part1(lines []string) int {
	cost := 0
	for _, l := range lines {
		err := CheckSyntax(l)
		if err == nil {
			continue
		}

		cost += syntaxErrorCost[err.got]
	}
	return cost
}

func part2(lines []string) int {
	costs := make([]int, 0)
	for _, l := range lines {
		if err := CheckSyntax(l); err != nil {
			continue
		}
		missing := MissingClosingCharacters(l)
		costs = append(costs, missingCost(missing))
	}

	sort.IntSlice(costs).Sort()

	return costs[(len(costs)-1)/2]
}

func missingCost(missing string) int {
	cost := 0
	for _, c := range missing {
		cost = 5*cost + completionCost[c]
	}
	return cost
}

type SyntaxError struct {
	input string
	want  rune
	got   rune
	idx   int
}

func (se SyntaxError) Error() string {
	return fmt.Sprintf("syntax error on %q: got %q, want %q at index %d", se.input, se.got, se.want, se.idx)
}

func CheckSyntax(line string) *SyntaxError {
	s := NewStack()
	for i, c := range line {
		if IsLeftBracket(c) {
			s.Push(c)
			continue
		}

		left, err := s.Pop()
		if err != nil {
			return &SyntaxError{input: line, got: c, want: '?', idx: i}
		}

		if match[left] != c {
			return &SyntaxError{input: line, got: c, want: match[left], idx: i}
		}
	}

	return nil
}

func MissingClosingCharacters(line string) string {
	s := NewStack()
	for _, c := range line {
		if IsLeftBracket(c) {
			s.Push(c)
			continue
		}

		_, _ = s.Pop()
	}

	var builder strings.Builder
	for !s.IsEmpty() {
		c, _ := s.Pop()
		builder.WriteRune(match[c])
	}

	return builder.String()
}

type Stack struct {
	data []rune
}

func NewStack() Stack {
	return Stack{data: make([]rune, 0)}
}

func (s *Stack) Push(value rune) {
	s.data = append(s.data, value)
}

func (s Stack) IsEmpty() bool {
	return len(s.data) == 0
}

var errorEmptyStack = errors.New("stack is empty")

func (s *Stack) Pop() (rune, error) {
	if len(s.data) == 0 {
		return 0, errorEmptyStack
	}
	res := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return res, nil
}

func IsLeftBracket(r rune) bool {
	_, ok := match[r]
	return ok
}

var match = map[rune]rune{
	'[': ']',
	'{': '}',
	'(': ')',
	'<': '>',
}

var syntaxErrorCost = map[rune]int{
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137,
}
var completionCost = map[rune]int{
	')': 1,
	']': 2,
	'}': 3,
	'>': 4,
}

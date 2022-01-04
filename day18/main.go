package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gverger/advent2021/utils"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {
	numbers := parseNumbers(lines)

	fmt.Println("Magnitude of the whole sum:", Part1(numbers))
	fmt.Println("Largest Magnitude", Part2(numbers))

	return nil
}

func Part1(numbers []FlatNumber) int {
	res := numbers[0]
	for _, n := range numbers[1:] {
		res = AddArray(res, n)
		res = ReduceArray(res)
	}

	return res.Magnitude()
}

func Part2(numbers []FlatNumber) int {
	max := 0

	for i, n1 := range numbers {
		for j, n2 := range numbers {
			if i == j {
				continue
			}
			mag := ReduceArray(AddArray(n1, n2)).Magnitude()

			if mag > max {
				max = mag
			}
		}
	}

	return max
}

func parseNumbers(lines []string) []FlatNumber {
	res := make([]FlatNumber, 0, len(lines))
	for _, l := range lines {
		res = append(res, ParseArray(l))
	}

	return res
}

const (
	Open  = -1
	Close = -2
)

type FlatNumber []int

func IsValue(n int) bool {
	return n >= 0
}

func ParseArray(line string) FlatNumber {
	res := make([]int, 0)

	current := 0
	parsingValue := false
	for _, c := range line {
		switch c {
		case '[':
			res = append(res, Open)
		case ']':
			if parsingValue {
				res = append(res, current)
				current = 0
				parsingValue = false
			}
			res = append(res, Close)
		case ',':
			if parsingValue {
				res = append(res, current)
				current = 0
				parsingValue = false
			}
		default:
			parsingValue = true
			n, e := strconv.Atoi(string(c))
			if e != nil {
				panic(e)
			}
			current = current*10 + n
		}
	}

	return res
}

func (n FlatNumber) Magnitude() int {
	multiplier := 1
	magnitude := 0

	needComma := false
	for _, p := range n {
		switch p {
		case Open:
			multiplier = multiplier * 3
			if needComma {
				multiplier = multiplier / 3
				multiplier = multiplier * 2
				needComma = false
			}
		case Close:
			multiplier = multiplier / 2
			needComma = true
		default:
			if needComma {
				multiplier = multiplier / 3
				multiplier = multiplier * 2
			}
			magnitude += multiplier * p
			needComma = true
		}
	}

	return magnitude
}

func (n FlatNumber) String() string {
	var b strings.Builder

	needComma := false
	for _, p := range n {
		switch p {
		case Open:
			if needComma {
				b.WriteRune(',')
				needComma = false
			}
			b.WriteRune('[')
		case Close:
			b.WriteRune(']')
			needComma = true
		default:
			if needComma {
				b.WriteRune(',')
			}
			b.WriteString(fmt.Sprintf("%d", p))
			needComma = true
		}
	}

	return b.String()
}

func AddArray(n1, n2 FlatNumber) FlatNumber {
	res := make(FlatNumber, 0)

	res = append(res, Open)
	res = append(res, n1...)
	res = append(res, n2...)
	res = append(res, Close)

	return res
}

func ReduceArray(n FlatNumber) FlatNumber {
	res, changed := Explode(n)

	for changed {
		for changed {
			res, changed = Explode(res)
		}
		res, changed = Split(res)
	}

	return res
}

func Explode(n FlatNumber) (FlatNumber, bool) {
	nestedPairIdx := n.FirstAtDepth(5)
	if nestedPairIdx == -1 {
		return n, false
	}

	leftNumber := n[nestedPairIdx+1]
	rightNumber := n[nestedPairIdx+2]

	for i := nestedPairIdx - 1; i >= 0; i-- {
		if IsValue(n[i]) {
			n[i] += leftNumber
			break
		}
	}
	for i := nestedPairIdx + 4; i < len(n); i++ {
		if IsValue(n[i]) {
			n[i] += rightNumber
			break
		}
	}

	res := make(FlatNumber, len(n)-3)
	copy(res, n[:nestedPairIdx])
	copy(res[nestedPairIdx+1:], n[nestedPairIdx+4:])

	return res, true
}

func Split(n FlatNumber) (FlatNumber, bool) {
	idx := n.FirstValueTooHigh()
	if idx == -1 {
		return n, false
	}

	res := make(FlatNumber, len(n)+3)
	copy(res, n[:idx])
	copy(res[idx+4:], n[idx+1:])
	res[idx] = Open
	res[idx+1] = n[idx] / 2
	res[idx+2] = (n[idx] + 1) / 2
	res[idx+3] = Close

	return res, true
}

func (n FlatNumber) FirstAtDepth(d int) int {
	depth := 0

	for i, c := range n {
		switch c {
		case Open:
			depth++
			if depth == d {
				return i
			}
		case Close:
			depth--
		}
	}

	return -1
}

func (n FlatNumber) FirstValueTooHigh() int {
	for i, c := range n {
		if !IsValue(c) {
			continue
		}
		if c >= 10 {
			return i
		}
	}

	return -1
}

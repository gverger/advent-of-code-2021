package main

import (
	"fmt"
	"strings"

	"github.com/gverger/advent2021/utils"
)

var Segments = map[byte][]int{
	'e': {0, 2, 6, 8},
	'b': {0, 4, 5, 6, 8, 9},
	'd': {2, 3, 4, 5, 6, 8, 9},
	'g': {0, 2, 3, 5, 6, 8, 9},
	'a': {0, 2, 3, 5, 6, 7, 8, 9},
	'c': {0, 1, 2, 3, 4, 7, 8, 9},
	'f': {0, 1, 3, 4, 5, 6, 7, 8, 9},
}

var Ns = map[int][]byte{
	1: {'c', 'f'},
	7: {'a', 'c', 'f'},
	4: {'b', 'c', 'd', 'f'},
	2: {'a', 'c', 'd', 'e', 'g'},
	3: {'a', 'c', 'd', 'f', 'g'},
	5: {'a', 'b', 'd', 'f', 'g'},
	0: {'a', 'b', 'c', 'e', 'f', 'g'},
	6: {'a', 'b', 'd', 'e', 'f', 'g'},
	9: {'a', 'b', 'c', 'd', 'f', 'g'},
	8: {'a', 'b', 'c', 'd', 'e', 'f', 'g'},
}

func main() {
	utils.Main(run)
}

func run(lines []string) error {

	fmt.Println("Part 1 =", part1(lines))
	fmt.Println("Part 2 =", part2(lines))
	return nil
}

func decodeInput(input string) Problem {
	p := Problem{}

	parts := strings.Split(input, "|")
	p.numbers = NumbersSlice(strings.Fields(parts[0]))
	p.output = NumbersSlice(strings.Fields(parts[1]))

	return p
}

func part1(lines []string) int {
	n := 0
	for _, l := range lines {
		problem := decodeInput(l)
		n += numberOfKnownNumbers(problem.output)
	}

	return n
}

func part2(lines []string) int {
	sum := 0
	for _, l := range lines {
		problem := decodeInput(l)
		sum += problem.associations()
	}

	return sum
}

func numberOfKnownNumbers(numbers []Number) int {
	res := 0

	for _, n := range numbers {
		switch len(n) {
		case 2, 3, 4, 7:
			res += 1
		}
	}

	return res
}

type Problem struct {
	numbers []Number
	output  []Number
}

type NumberValue struct {
	numberToValue map[Number]int
	valueToNumber map[int]Number
}

func NewNumberValue() NumberValue {
	return NumberValue{numberToValue: make(map[Number]int), valueToNumber: make(map[int]Number)}
}

func (nv *NumberValue) Associate(number Number, value int) {
	nv.numberToValue[number] = value
	nv.valueToNumber[value] = number
}

func (nv NumberValue) NumberForValue(value int) Number {
	return nv.valueToNumber[value]
}

func (nv NumberValue) ValueForNumber(number Number) int {
	return nv.numberToValue[number]
}

type WireSegment struct {
	wireToSegment map[rune]rune
	segmentToWire map[rune]rune
}

func NewWireSegment() WireSegment {
	return WireSegment{wireToSegment: make(map[rune]rune), segmentToWire: make(map[rune]rune)}
}

func (ws *WireSegment) Associate(wire, segment rune) {
	ws.wireToSegment[wire] = segment
	ws.segmentToWire[segment] = wire
}

func (ws WireSegment) SegmentForWire(wire rune) rune {
	return ws.wireToSegment[wire]
}

func (ws WireSegment) WireForSegment(segment rune) rune {
	return ws.segmentToWire[segment]
}

func (p Problem) associations() int {
	nv := NewNumberValue()
	ws := NewWireSegment()

	for _, n := range p.numbers {
		switch len(n) {
		case 2:
			nv.Associate(n, 1)
		case 3:
			nv.Associate(n, 7)
		case 4:
			nv.Associate(n, 4)
		case 7:
			nv.Associate(n, 8)
		}
	}

	wires := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'}
	wireNumber := make(map[rune][]Number)
	for _, w := range wires {
		wireNumber[w] = make([]Number, 0)
	}

	for _, n := range p.numbers {
		for w := range n.Wires() {
			wireNumber[w] = append(wireNumber[w], n)
		}
	}

	for w, n := range wireNumber {
		switch len(n) {
		case 4:
			ws.Associate(w, 'e')
		case 6:
			ws.Associate(w, 'b')
		case 9:
			ws.Associate(w, 'f')
		}
	}

	// Find which wire correspond to segment c
	// It is the unknown wire from number 1, since we already now the segment f
	for w := range nv.NumberForValue(1).Wires() {
		if ws.SegmentForWire(w) == 'f' {
			continue
		}
		ws.Associate(w, 'c')
	}

	// 2, 3 and 5 have 5 segments.
	// 5 is the only one with segment b
	// 3 has segment f on
	for _, n := range p.numbers {
		if len(n) != 5 {
			continue
		}

		if n.HasWire(ws.WireForSegment('b')) {
			nv.Associate(n, 5)
		} else if n.HasWire(ws.WireForSegment('f')) {
			nv.Associate(n, 3)
		} else {
			nv.Associate(n, 2)
		}
	}

	// 0, 6 and 9 have 6 segments
	// 6 is the only one without segment c
	// 0 has segment e on
	for _, n := range p.numbers {
		if len(n) != 6 {
			continue
		}

		if !n.HasWire(ws.WireForSegment('c')) {
			nv.Associate(n, 6)
		} else if n.HasWire(ws.WireForSegment('e')) {
			nv.Associate(n, 0)
		} else {
			nv.Associate(n, 9)
		}
	}

	res := 0
	for _, o := range p.output {
		for n := range nv.numberToValue {
			if o.Equal(n) {
				res *= 10
				res += nv.ValueForNumber(n)
				break
			}
		}

	}

	return res
}

type Number string

func NumbersSlice(numbers []string) []Number {
	res := make([]Number, len(numbers))

	for i, n := range numbers {
		res[i] = Number(n)
	}

	return res
}

func (n Number) Equal(other Number) bool {
	if len(n) != len(other) {
		return false
	}

	chars := n.Wires()

	for _, c := range other {
		if !chars[c] {
			return false
		}
	}

	return true
}

func (n Number) Wires() map[rune]bool {
	wires := make(map[rune]bool)

	for _, c := range n {
		wires[c] = true
	}

	return wires
}

func NumbersFilter(numbers []Number, filter func(n Number) bool) []Number {
	res := make([]Number, 0)

	for _, n := range numbers {
		if filter(n) {
			res = append(res, n)
		}
	}

	return res
}

func (n Number) HasWire(wire rune) bool {
	return n.Wires()[wire]
}

func NumbersWithWire(numbers []Number, wire rune) []Number {
	return NumbersFilter(numbers, func(n Number) bool { return n.HasWire(wire) })

}

func NumbersWithoutWire(numbers []Number, wire rune) []Number {
	return NumbersFilter(numbers, func(n Number) bool { return !n.HasWire(wire) })
}

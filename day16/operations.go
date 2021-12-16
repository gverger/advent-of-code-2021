package main

import "github.com/gverger/advent2021/utils"

const (
	TypeSum     = 0
	TypeProduct = 1
	TypeMinimum = 2
	TypeMaximum = 3
	TypeGreater = 5
	TypeLess    = 6
	TypeEqual   = 7
)

type operation func(packets []Packet) int

var operationForType = map[int]operation{
	TypeSum:     sum,
	TypeProduct: product,
	TypeMinimum: minimum,
	TypeMaximum: maximum,
	TypeLess:    less,
	TypeGreater: greater,
	TypeEqual:   equal,
}

func sum(packets []Packet) int {
	return reduce(packets, func(i1, i2 int) int { return i1 + i2 })
}

func product(packets []Packet) int {
	return reduce(packets, func(i1, i2 int) int { return i1 * i2 })
}

func minimum(packets []Packet) int {
	return reduce(packets, func(i1, i2 int) int { return utils.Min(i1, i2) })
}

func maximum(packets []Packet) int {
	return reduce(packets, func(i1, i2 int) int { return utils.Max(i1, i2) })
}

func reduce(packets []Packet, combine func(int, int) int) int {
	res := packets[0].Compute()
	for _, sub := range packets[1:] {
		res = combine(res, sub.Compute())
	}
	return res
}

func less(packets []Packet) int {
	return comparison(packets, func(i1, i2 int) bool { return i1 < i2 })
}

func greater(packets []Packet) int {
	return comparison(packets, func(i1, i2 int) bool { return i1 > i2 })
}

func equal(packets []Packet) int {
	return comparison(packets, func(i1, i2 int) bool { return i1 == i2 })
}

func comparison(packets []Packet, compare func(int, int) bool) int {
	a := packets[0].Compute()
	b := packets[1].Compute()
	if compare(a, b) {
		return 1
	}
	return 0
}

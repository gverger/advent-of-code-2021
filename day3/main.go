package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gverger/advent2021/utils"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {
	values, err := StringSlice(lines).mapToCounts()
	if err != nil {
		return fmt.Errorf("cannot convert to counts: %w", err)
	}

	fmt.Printf("%#v: %v/%v\n", values, values.Gamma(), values.Epsilon())

	fmt.Println("Gamma x Epsilon:", values.Gamma()*values.Epsilon())
	fmt.Printf("Oxygen : %d, CO2 : %d\n", StringSlice(lines).Oxygen(), StringSlice(lines).CO2())
	fmt.Println("Oxygen x CO2: ", StringSlice(lines).Oxygen()*StringSlice(lines).CO2())

	return nil
}

type Counts []int

func (c Counts) Gamma() int {
	res := 0
	for _, b := range c {
		res *= 2
		if b > 0 {
			res += 1
		}
	}

	return res
}

func (c Counts) Epsilon() int {
	res := 0
	for _, b := range c {
		res *= 2
		if b < 0 {
			res += 1
		}
	}

	return res
}

type StringSlice []string

func (data StringSlice) mapToCounts() (Counts, error) {
	counts := make([]int, 0)
	for _, c := range data[0] {
		switch c {
		case '1':
			counts = append(counts, 1)
		case '0':
			counts = append(counts, -1)
		default:
			return nil, errors.New("unauthorized character")
		}
	}

	for _, text := range data[1:] {
		for i, c := range text {
			switch c {
			case '1':
				counts[i] += 1
			case '0':
				counts[i] -= 1
			default:
				return nil, errors.New("unauthorized character")
			}
		}

	}

	return counts, nil
}

func (data StringSlice) Filter(sel func(string) bool) StringSlice {
	res := make(StringSlice, 0)

	for _, value := range data {
		if sel(value) {
			res = append(res, value)
		}
	}

	return res
}

func (data StringSlice) Oxygen() int {
	d := make(StringSlice, len(data))
	copy(d, data)
	nb := len(d[0])

	for i := 0; i < nb; i++ {
		counts, _ := d.mapToCounts()
		c := counts[i]
		filter := '1'
		if c < 0 {
			filter = '0'
		}
		d = d.Filter(func(s string) bool {
			return s[i] == byte(filter)
		})
		if len(d) == 1 {
			res, _ := strconv.ParseInt(d[0], 2, 0)
			return int(res)
		}
	}

	return -1
}

func (data StringSlice) CO2() int {
	d := make(StringSlice, len(data))
	copy(d, data)
	nb := len(d[0])

	for i := 0; i < nb; i++ {
		counts, _ := d.mapToCounts()
		c := counts[i]
		filter := '1'
		if c >= 0 {
			filter = '0'
		}
		d = d.Filter(func(s string) bool {
			return s[i] == byte(filter)
		})
		if len(d) == 1 {
			res, _ := strconv.ParseInt(d[0], 2, 0)
			return int(res)
		}
	}

	return -1
}

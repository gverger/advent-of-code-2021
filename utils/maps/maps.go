package maps

import (
	"strconv"
	"strings"
)

type Strings []string

func (data Strings) ToInts() ([]int, error) {
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

type Ints []int

func (data Ints) ToStrings() ([]string, error) {
	values := make([]string, 0, len(data))
	for _, text := range data {
		current := strconv.Itoa(text)
		values = append(values, current)
	}

	return values, nil
}

type String string

func (data String) ToInts() ([]int, error) {
	values := make([]int, 0, len(data))
	for _, text := range data {
		current, err := strconv.Atoi(string(text))
		if err != nil {
			return nil, err
		}

		values = append(values, current)
	}

	return values, nil
}

// "12345".ChunkEvery(2) returns -> [12 34], leaving the 5 out
func (data String) ChunkEvery(chunkSize int) []string {
	res := make([]string, 0)
	var builder strings.Builder
	for i, r := range data {
		builder.WriteRune(r)
		if i > 0 && (i+1)%chunkSize == 0 {
			res = append(res, builder.String())
			builder.Reset()
		}
	}

	return res
}

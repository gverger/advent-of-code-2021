package maps

import "strconv"

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

package filters

type Strings []string

func (data Strings) KeepIf(condition func(string) bool) []string {
	res := make([]string, 0)

	for _, value := range data {
		if condition(value) {
			res = append(res, value)
		}
	}

	return res
}

func (data Strings) Uniq() []string {
	set := make(map[string]bool)
	for _, e := range data {
		set[e] = true
	}

	res := make([]string, 0, len(set))

	for e := range set {
		res = append(res, e)
	}

	return res
}

type Ints []int

func (data Ints) Uniq() Ints {
	set := make(map[int]bool)
	for _, e := range data {
		set[e] = true
	}

	res := make([]int, 0, len(set))

	for e := range set {
		res = append(res, e)
	}

	return res
}

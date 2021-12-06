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

package debug

import "fmt"

type Trace struct {
	prefix string
	on     bool
}

func (t Trace) printf(format string, args ...interface{}) {
	if t.on {
		fmt.Printf(t.prefix+format+"\n", args...)
	}
}

func (t Trace) println(args ...interface{}) {
	if t.on {
		fmt.Print(t.prefix)
		fmt.Println(args...)
	}
}

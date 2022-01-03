package debug

import "fmt"

type Trace struct {
	Prefix string
	On     bool
}

func (t Trace) Printf(format string, args ...interface{}) {
	if t.On {
		fmt.Printf(t.Prefix+format, args...)
	}
}

func (t Trace) Printfln(format string, args ...interface{}) {
	if t.On {
		fmt.Printf(t.Prefix+format+"\n", args...)
	}
}

func (t Trace) Println(args ...interface{}) {
	if t.On {
		fmt.Print(t.Prefix)
		fmt.Println(args...)
	}
}

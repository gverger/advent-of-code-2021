package main

import "fmt"

func (p LiteralPacket) String() string {
	return fmt.Sprintf("LIT(%v, value=%d)", p.Header, p.value)
}

func (h Header) String() string {
	return fmt.Sprintf("H[v%d t%d]", h.Version(), h.TypeID())
}

func (p OperatorPacket) String() string {
	packets := make([]string, len(p.packets))

	for i, sub := range p.packets {
		packets[i] = fmt.Sprint(sub)
	}

	return fmt.Sprintf("%s(%v, packets=%v)", stringForType[p.TypeID()], p.Header, packets)
}

var stringForType = map[int]string{
	TypeSum:     "SUM",
	TypeProduct: "PROD",
	TypeMinimum: "MIN",
	TypeMaximum: "MAX",
	TypeLess:    "LESS",
	TypeGreater: "MORE",
	TypeEqual:   "EQ",
}

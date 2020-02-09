package main

import (
	"fmt"
	"unsafe"
)

func newString(s string) *string {
	return &s
}

func main() {
	m1 := map[string]*string{
		"a": newString("a"),
		"b": newString("b"),
		"c": newString("c"),
	}

	m2 := m1

	fmt.Printf("m1 addr = %v\n", unsafe.Pointer(&m1))
	fmt.Printf("m2 addr = %v\n", unsafe.Pointer(&m2))
}

package main

import (
	"fmt"
	"reflect"
)

func main() {
	// a is an array of type [4]uint
	a := [4]byte{'h', 'a', 'l', 'o'}
	fmt.Printf("a is an array of type %T, a=%v\n", a, a)

	// b is an array of type [5]uint
	b := [5]byte{'h', 'e', 'l', 'l', 'o'}
	fmt.Printf("b is an array of type %T, b=%v\n", b, b)

	// so, a and b are not of the same type
	fmt.Printf("are a and b of the same type? %v\n", reflect.TypeOf(a)==reflect.TypeOf(b))

	// c is a slice with 4 elements
	c := []byte{'h', 'a', 'l', 'o'}
	fmt.Printf("c is a slice of type %T, c=%v\n", c, c)

	// we can also make an zero value slice, slice is nil
	var d []byte
	fmt.Printf("Is d nil? %v\n", d == nil)

	// slice is a []T, as we can see that both len and cap take 'v Type' as parameter
	fmt.Printf("len of d: %d, cap of d:%d\n", len(d), cap(d))

	// we can't expand d because d is a slice, thus the following does not work
	// d[0] = 'a'

	// a sliced slice is still a slice, but the new slice is still referencing the same storage of the original slice
	e := []byte{1, 2, 3, 4, 5}
	fmt.Printf("e=%v\n", e)

	// f = [3,4]
	f := e[2:4]

	// we modify the element (values) of f, f=[120, 120]
	f[0] = 'x'
	f[1] = 'x'

	// since f and e shares the same storage, so e=[1, 2, 120, 120, 5]
	fmt.Printf("after modification, e=%v\n", e)

	// copy solves the 'share-storage' problem
	e = []byte{1, 2, 3, 4, 5}
	g := make([]byte, 2)
	copy(g, e[2:4])
	g[0] = 'x'
	g[1] = 'x'
	fmt.Printf("after modification [with copy treatment], e=%v\n", e)

	// append solves the capacity expanding limit problem
	e = append(e, 'x')
	fmt.Printf("after append, e=%v\n", e)

}



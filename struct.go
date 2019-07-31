package main

import "unsafe"

func main() {
	var sn1 struct {n int8}
	var sn2 struct {n int16}
	var sn3 struct {n int32}
	var sn4 struct {n int64}
	println(unsafe.Sizeof(sn1), unsafe.Sizeof(sn2), unsafe.Sizeof(sn3), unsafe.Sizeof(sn4))
	// 1, 2, 4, 8 (same for uint)

	var sf1 struct {f float32}
	var sf2 struct {f float64}
	println(unsafe.Sizeof(sf1), unsafe.Sizeof(sf2))
	// 4, 8

	var si struct {i interface{}}
	println(unsafe.Sizeof(si))
	// 16

	var s1, s2 struct{}
	println(unsafe.Sizeof(s1) == unsafe.Sizeof(s2),  unsafe.Sizeof(s2) == 0)
	// true, true
	// empty struct does not have any allocation, and the size is zero

	println(&s1, &s2)
	// empty structs share the same address
	// NOTE! you cant compare the address using '==', thus, &s1 != &s2, but they are the same address.

	var s3, s4 struct { x struct{} }
	println(unsafe.Sizeof(s3) == unsafe.Sizeof(s4),  unsafe.Sizeof(s3) == 0)
	// true, true
	// a struct that only has an empty struct member is still considered an empty struct

	a := make([]struct{}, 1)
	b := make([]struct{}, 9999)

	println(&a, &b)
	// two empty struct slices does not have the same address (because they are no longer empty)
	//  but their backing value (array) are the same

	println(unsafe.Sizeof(a) == unsafe.Sizeof(b), unsafe.Sizeof(a))
	// true, true
	// an empty slice has the allocation of its header size (12 in 32 bit and 24 in 64 bit machine) no matter
	// what type of slice it is

}

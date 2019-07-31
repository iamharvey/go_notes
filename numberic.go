package main

import "math"

func main() {

	println("uint8", "min:", uint8(0), "max:", ^uint8(0))
	println("uint16", "min:", uint16(0), "max:", ^uint16(0))
	println("uint32", "min:", uint32(0), "max:", ^uint32(0))
	println("uint64", "min:", uint64(0), "max:", ^uint64(0))
	// for unsigned numbers, it is 0 - (2^n - 1), where n \in [8, 16, 32, 64]

	println("int8", "min:", -(int(^uint8(0) >> 1)) - 1, "max:", int(^uint8(0) >> 1))
	println("int16", "min:", -(int(^uint16(0) >> 1)) - 1, "max:", int(^uint16(0) >> 1))
	println("int32", "min:", -(int(^uint32(0) >> 1)) - 1, "max:", int(^uint32(0) >> 1))
	println("int64", "min:", -(int(^uint64(0) >> 1)) - 1, "max:", int(^uint64(0) >> 1))
	// for signed ones, the range the symmetric w.r.t zero
	// we can also use `math` package to get the MIN and MAX, e.g., math.MinInt8, math.MaxInt32

	// the set of all IEEE-754 320bit floating-point numbers
	println("float32", "max", math.MaxFloat32)
	println("float64", "max", math.MaxFloat64)

	// the set of all complex numbers with float real and imaginary parts
	println("complex 64", "max", complex(math.MaxFloat32, math.MaxFloat32))
	println("complex 128", "max", complex(math.MaxFloat64, math.MaxFloat64))

	// byte is alias for uint8, so the min max are the same as uint8
	// rune is the alias for int32, so the min max are the same as int32
}

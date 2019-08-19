package main

type A struct {
	x int32
}

type B struct {
	a *A
}

type C struct {
	b *B
}

func (b *B) NewA() *A {
	if b != nil {
		return b.a
	}
	return nil
}

func (c *C) NewB() *B {
	if c != nil {
		return c.b
	}

	return nil
}

func main() {
	// this will not compile, we can not assign nil without explicit type
	// because Go has no idea what it is: a pointer? a string, an array?
	//a := nil
	//println(a)

	var a *int = nil
	var b *int = nil

	println(a, b, &a, &b)

	println(a == b, a == nil)
	// true, true

	var c *float32 = nil
	println(c, &c)
	// this wont compile, as we can not compare nil of different type
	// println(c == a)

	// in Go, an empty struct permits to have call on its method,
	// and it wont panic (throw null pointer error like java)
	y := C{}
	println(y.NewB())
	// 0x0
	println(y.NewB().NewA())
	// 0x0

}

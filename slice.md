# About Slice

Slice type is an abstraction on top of array type. 

## Essentials about array
- An array definition specifies its length and element type, e.g., [4]byte is an array of type [4]uint with 4 elements. 
- two arrays of different size have different type, e.g., [4]byte is an array of type [4]uint while [5]byte is an array of type [5]uint.
- Array size is fixed.
- Arrays are values. When you are assigning or passing around an array value, you are actually copying the value of the array.

## Essentials about slice
- Slice is []T (see `len()` and `cap()`, the argument is of type []T).
- Defining a slice require no size, e.g., var b []byte.
- We can use `make()` to define a slice.
- An zero-value slice is `nil`, and with the size of zero and capacity of zero. 
  ```
  var b []byte
  fmt.Print(b == nil) // true
  fmt.Print(len(b), cap(b)) // 0, 0
  ```
- When make a slice with size, the initial value is assigned:
    - for byte, int, float, the initial value is zero;
    - for string, the initial value is empty string "";
    - for struct, the initial value is empty struct {};
    - for interface, the initial value is nil.
- Slices can be sliced using `:`, this is similar with that in Python programming language, except that the index can not be negative.
- A slice is a descriptor of an array segment. It contains: pointer (to the array), segment length, and capacity.
- When we are doing slicing, we are actually not copying the values 'out', but creating a new slice, but referencing to the sliced value(s)
  in the original slice. For e.g.,
  ```
	e := []byte{1, 2, 3, 4, 5}
	// e = [1, 2, 3, 4, 5]

	f := e[2:4]
    // f = [3,4]

	f[0] = 'x'
	f[1] = 'x'
    // we modify the elements (values) of f, thus, f=[120, 120]
	// now, e=[1, 2, 120, 120, 5]
  ```
- To avoid 'referencing to original slice', we can use copy
- To expand the capacity of a slice, we can use append

## Reference:
[1. Go Slices: usage and internals.](https://blog.golang.org/go-slices-usage-and-internals)
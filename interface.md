# About Interface

## The Basics
Interface is one of the essential in Golang. It is a built-in type that is defined as a set of method signatures.
A value of interface can hold any value that implements those methods [1].

E.g.,
```
type Robot interface {
	Forward(speed, dur int) int
	Backward(speed, dur int) int
	TurnLeft(degree int) int
	TurnRight(degree int) int
}

type MyRobot struct {
	Robot
}

func (bot *MyRobot) Forward(speed, dur int) (meter int) {
	d := speed * dur
	fmt.Printf("bot has moved %d meters forward ... \n", d)
	return d
}

func (bot *MyRobot) Backward(speed, dur int) (meter int) {
	d := speed * dur
	fmt.Printf("bot has moved %d meters backward ... \n", d)
	return d
}

func (bot *MyRobot) TurnLeft(deg int) int {
	fmt.Printf("bot has turnd %d degree left ... \n", deg)
	return deg
}

func main() {
	var bot MyRobot
	bot.Forward(1, 2)
	bot.TurnLeft(30)
	bot.Backward(1, 2)
	bot.TurnRight(45)
}

```

The above example does not implement `TurnRight` method. If we run the program, error occurs:
```
$ go run interface.go
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x30 pc=0x109ae8c]
```

Now we add `TurnRight` function:
```
func (bot *MyRobot) TurnRight(deg int) int {
	fmt.Printf("bot has turnd %d degree left ... \n", deg)
	return deg
}
```

Run the program again, the robot now 'is functioning correctly'.
```
bot has moved 2 meters forward ... 
bot has turnd 30 degree left ... 
bot has moved 2 meters backward ... 
bot has turnd 45 degree left ... 
```

We can implicitly declare the interface, e.g.,
```
...
type MyRobot2 struct {}

func (bot *MyRobot2) Forward(speed, dur int) (meter int) {
	d := speed * dur
	fmt.Printf("bot has moved %d meters forward ... \n", d)
	return d
}

func (bot *MyRobot2) Backward(speed, dur int) (meter int) {
	d := speed * dur
	fmt.Printf("bot has moved %d meters backward ... \n", d)
	return d
}

func (bot *MyRobot2) TurnLeft(deg int) int {
	fmt.Printf("bot has turnd %d degree left ... \n", deg)
	return deg
}

func (bot *MyRobot2) TurnRight(deg int) int {
	fmt.Printf("bot has turnd %d degree left ... \n", deg)
	return deg
}

...

func main() { 
    ...
    fmt.Println("\nstarting the second robot ... ")
    var bot2 Robot = &MyRobot2{}
    bot2.Forward(1, 3)
    bot2.TurnLeft(25)
    bot2.Backward(2, 2)
    bot2.TurnRight(30)
}
```

It works the same as the previous approach:
```
$ go run iterface.go
bot has moved 2 meters forward ... 
bot has turnd 30 degree left ... 
bot has moved 2 meters backward ... 
bot has turnd 45 degree left ... 

starting the second robot ... 
bot has moved 3 meters forward ... 
bot has turnd 25 degree left ... 
bot has moved 4 meters backward ... 
bot has turnd 30 degree left ... 
```

We can also make a single type (e.g., MyRobot) implements multiple interfaces, e.g.,
```
fmt.Println("\nstarting two agents, one on the ground, one in the sky ... ")
var agent1 Robot = &MyRobot{}
var agent2 Flybot = &MyRobot{}
agent1.Forward(1, 2)
agent1.TurnLeft(30)
agent2.Backward(1, 2)
agent2.TurnRight(45)
```

## Interface Embedding
An interface can also hold another interface, e.g.,
```
type I interface {
    i()
}

type J interface {
   I
   j()
}

type K interface {
   J
   k()
}
```
In the above example, all the method from I, J will be added to K. Thus, we need implement all the function `i()`,
`j()` Nd `k()`, otherwise run-time error occurs. 

Two rules are not allowed to violate:
- **circular embedding**. E.g., I holds J, while J holds I.
- **duplicate function**. E.g., I holds `i()` where J also holds `i()`. Note! The name matters not the signature.

## Nil Interface Is Not Empty
A nil interface holds neither value nor concrete type. Thus, calling a method on a nil interface 
can have a run-time error. E.g.,
```
var bot3 Robot
fmt.Printf("%v, %T", bot3, bot3)
bot3.TurnRight(30)
```
The above code throws a run-time error says that the interface is a nil pointer dereference.

Note that an empty interface (`interface{}`) is not a nil interface. It is a valid interface that holds no 
concrete values, e.g., 
```
var i interface{}
fmt.Printf("%v, %T\n", i, i)

i = 123
_, ok := i.(int)
fmt.Printf("%v, %T, %v\n", i, i, ok==true)

i = "hello interface"
_, ok = i.(string)
fmt.Printf("%v, %T, %v\n", i, i, ok==true)
```

The output is:
```
& go run interface.go
...
<nil>, <nil>
123, int, true
hello interface, string, true
```

As it is shown above, we can also use `.(<type>)` to check out the interface's value's concrete type.

When we play with the pointer of a type within the same function, such type is indeed a nil value type. E.g., 
```
func main() {
...
    var b4 *MyRobot
    fmt.Printf("%v, %t, %v", b4, b4, b4 == nil)
    // the result is: <nil>, *main.MyRobot, true
...
}
```

But the truth can not be hold if such pointer returns from a function without initialisation. E.g.,
```
func getBot() Robot {
	var b4 *MyRobot
	return b4
}
...
func main() {
    ...
    fmt.Printf("%v, %t, %v", getBot(), getBot(), getBot() == nil)
    // the result is: <nil>, *main.MyRobot, false
    ...
}
```

## Re-assign
An interface variable's value can be re-assigned to another interface variable, and the function does not 
have to be equally the same. But the re-assignment target interface should hold the same function signature. 
E.g. the following example throws a run-time error says that `j()` is not implemented in I.
```
type I interface {
	i()
}

type I2 interface {
	j()
}

type T struct{}

func (T) i() {}
func (T) j() {}

func main() {
	var v1 I1 = T{}
	var v2 I2 = v1
	_ = v2
}

```

Adding `j()` to I solves the problem:
```
type I interface {
	i()
    j()
}
```

## Type Assertion
Note that, we can not explicitly concert an interface value to a certain type. E.g., the following example 
does not work:
```
type I interface {
	i()
}

type T struct{}

func (T) i() {}

func main() {
	var i I = T{}
    fmt.Println(T(i))
}

```

Instead, type assertion is needed:
```
...
func main() {
	var i I = T{}
    fmt.Println(i.(T))
}
...

```


When there are two different interface types, we can convert one to another type. E.g:
```

type I interface {
	i()
}

type J interface {
	j()
}

type T struct {
	name string
}

func (T) i() {}
func (T) j() {}

func main() {
	var i I = T{"hi"}
	var j J
	j, ok := i.(J)
	fmt.Printf("%T %v %v\n", j, j, ok)
}
```


But the following example throws an error:
```
type I interface {
    i()
}
type T1 struct{}
func (T1) i() {}
type T2 struct{}
func (T2) i() {}
func main() {
    var i1 I = T1{}
    i2 := i1.(T2)
    fmt.Printf("%T\n", i2)
}
```

This is because the dynamic type of `i1` does not match T2. The type assertion also fails when the re-assignment 
target is `nil`. E.g.,
```
package main

type I interface {
	i()
}

type T struct{}

func (T) i() {}

func main() {
	var i I // should initialise it first
    fmt.Println(i.(T))
}
```

We can use `.(type)` to dynamically check the type of an interface variable. This helps us to perform function that 
responses dynamically according to different type:
```
type I interface {
	i()
}

type T1 struct{}

func (T1) i() {}

type T2 struct{}

func (T2) i() {}

func main() {
	var i I = T1{}
	switch aux := 1; v.(type) {
	case nil:
		fmt.Println("nil")
	case T1:
		fmt.Println("T1", aux)
	case T2:
		fmt.Println("T2", aux)
	}
}

```

The above example matches the condition of T1. We can also check the pointer type, e.g.:
```
var i1 *T1
var i2 I = i1
switch t := v.(type) {
case nil:
         fmt.Println("nil")
case *T1, *T2:
         fmt.Printf("%T is nil: %v\n", t, t == nil)
}
```

The result matches the condition of *T1, *T2, because pointer that points to a `nil` is not `nil`.

The full source code can be found [here](interface.go)

## Reference:
1. [Go interface](https://tour.golang.org/methods/9)
2. [Interface In Go (Part II) by Michał Łowicki](https://medium.com/golangspec/interfaces-in-go-part-ii-d5057ffdb0a6)
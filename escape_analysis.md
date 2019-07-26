# About escape analysis

Escape analysis is an important task for GO compiler, it determines where a variable is 
constructed at compilation time (stack or heap)?

Escape analysis is regulated based on the escape analysis algorithm. It varies from language to language.
In short:

_If the compiler **CAN NOT** figure out if a value will be referenced after it returns from a function, 
in order to maintain the program's **integrity**, it moves the value to the heap. Otherwise, it will be constructed 
in the program's stack frame._

Let's see an example.
```
package main

type hero struct {
	name 		string
	superpower	string
}

//go:noinline
func createSuperMan() hero {
	h := hero {
		name: 		"Superman",
		superpower: "X-ray vision",
	}
	println("Superman ", &h)
	return h
}

//go:noinline
func createTheFlash() *hero {
	h := hero {
		name: 		"The Flash",
		superpower: "Super speed",
	}
	println("The Flash ", &h)
	return &h
}

/*
	This example use go:noinline (no space between // and go) to prevent the compiler from inlining the code
	for these functions directly in main. Inlining would erase the function calls and complicate the example.
*/
func main() {
	h1 := createSuperMan()

	// h will escape to heap, after createTheFlash returns its value, no matter if h2 is declared or not
	h2 := createTheFlash()

	println("Superman ", &h1, "The Flash", &h2)
}
```

The above example is a mimicry of the example code 
in [William Kennedy's post about escape analysis](https://www.goinggo.net/2017/05/language-mechanics-on-escape-analysis.html). 
I just replace the user struct with my hero struct. The `main()` function creates two heroes: Superman and The Flash.

When we create Superman (see `createSuperMan()`), we first create a hero struct `h` and return its value copy to the main function.
`h` stays in the stack frame. When it returns to the main function, the value of `h` is still kept in the stack but is no longer 
valid to access. The memory block will be re-framed and re-initialised after the next function call.

When we create The Flash (see `createTheFlash()`), we first create a hero struct `h` and return its reference (i.e. the address) to 
the main function. In this case, `h` is 'shared' with the main function. Thus, the compiler considers 'that sharing' can potentially 
hurt the program integrity and moves `h` to heap. If we run with gc flags `-gcflags "-m -m"`, we can see some information about variable escapes to 
heap:

```
$ go run -gcflags "-m -m" escape_analysis.go
./escape_analysis.go:9:6: cannot inline createSuperMan: marked go:noinline
./escape_analysis.go:19:6: cannot inline createTheFlash: marked go:noinline
./escape_analysis.go:32:6: cannot inline main: function too complex: cost 133 exceeds budget 80
./escape_analysis.go:14:23: createSuperMan &h does not escape
./escape_analysis.go:25:9: &h escapes to heap
./escape_analysis.go:25:9:      from ~r0 (return) at ./escape_analysis.go:25:2
./escape_analysis.go:20:2: moved to heap: h
./escape_analysis.go:24:24: createTheFlash &h does not escape
./escape_analysis.go:38:23: main &h1 does not escape
./escape_analysis.go:38:41: main &h2 does not escape
Superman  0xc000046710
The Flash  0xc00000c040
Superman  0xc000046768 The Flash 0xc000046760
```

Because `&h` escapes to heap, `h` thus will be constructed in the heap and both main function and the two hero creation 
functions can **ONLY INDIRECTLY ACCESS** to it.

**REMEMBER!** Escape analysis only happens at the compilation time, thus only GO compiler knows if a variable should be 
constructed in stack or heap. Although keeping variables sits in stack can avoid latency introduced by GC that operates on heap. 
We should not be obsessive about figuring out where a variable puts when we are coding. Instead, make sure our program is coded 
correctly (in terms of the business logic) and has good readability and interpretability should always be our primary goal.




